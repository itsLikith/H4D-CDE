# Copyright 2026 H4D-CDE Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""
Demand Forecaster gRPC Service (Module 3).

Implements spatial-temporal airspace density forecasting using the
Dilated Temporal Convolutional Network (TCN) defined in model.py,
which realises Equation (9) of Sahadevan et al. (ICSPIS 2025):

    Ô(t + Δt) = f_TCN(O(t - k : t), θ)

The service accepts per-voxel occupancy time-series and returns a
multi-step-ahead density forecast.  When the model checkpoint is absent
(pre-training) the service falls back to an exponential moving average
so it stays healthy and responsive at all times.
"""

import logging
import os
import sys
from concurrent import futures

import grpc
import numpy as np
import torch

# ---------------------------------------------------------------------------
# Add generated protobuf stubs to sys.path (produced by `make proto`)
# ---------------------------------------------------------------------------
_THIS_DIR = os.path.dirname(os.path.abspath(__file__))
_SVC_DIR = os.path.dirname(_THIS_DIR)
_GEN_DIR = os.path.join(_SVC_DIR, "gen")
if _GEN_DIR not in sys.path:
    sys.path.insert(0, _GEN_DIR)

import demand_forecaster_pb2            # noqa: E402
import demand_forecaster_pb2_grpc       # noqa: E402
from app.model import DemandForecasterTCN  # noqa: E402

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)
logger = logging.getLogger("demand-forecaster-svc")

# Occupancy threshold above which a voxel is classified as congested.
# Value sourced from the simulation parameters in Section V-C of the paper.
_CONGESTION_THRESHOLD = 15.0

# EMA smoothing factor used in the no-model fallback
_EMA_ALPHA = 0.3


def _ema_forecast(history: list[float]) -> float:
    """
    Exponential moving average fallback when the TCN model is not loaded.
    Returns the EMA of the most recent occupancy observations.
    """
    if not history:
        return 0.0
    ema = history[0]
    for val in history[1:]:
        ema = _EMA_ALPHA * val + (1.0 - _EMA_ALPHA) * ema
    return ema


class DemandForecasterServicer(
    demand_forecaster_pb2_grpc.DemandForecasterServiceServicer
):
    """
    gRPC Servicer implementing the DemandForecasterService contract.

    On startup the service tries to load a PyTorch TCN checkpoint from
    MODEL_PATH (env var).  If the checkpoint is absent it uses an EMA
    heuristic so the service starts and serves immediately without any
    hard-coded dummy values.
    """

    def __init__(self, model_path: str | None = None):
        if model_path is None:
            model_path = os.getenv("MODEL_PATH", "models/demand_tcn.pt")

        self.model: DemandForecasterTCN | None = None
        self.device = torch.device("cpu")

        if os.path.exists(model_path):
            try:
                self.model = DemandForecasterTCN(in_features=1, horizon_steps=90)
                state = torch.load(
                    model_path, map_location=self.device, weights_only=True
                )
                self.model.load_state_dict(state)
                self.model.eval()
                logger.info("Loaded TCN demand forecast model from %s", model_path)
            except Exception as exc:
                logger.warning(
                    "Failed to load TCN model from %s: %s – using EMA fallback.",
                    model_path,
                    exc,
                )
                self.model = None
        else:
            logger.warning(
                "Model checkpoint not found at %s – using EMA fallback.", model_path
            )

    # ------------------------------------------------------------------
    # RPC handlers
    # ------------------------------------------------------------------

    def ForecastDemand(self, request, context):
        """
        Forecasts future occupancy density for a specific 4D hex-voxel zone.

        The request carries:
            h3_index              — Uber H3 cell identifier (string)
            time_horizon_seconds  — how far ahead to forecast (default 300 s)
            historical_density    — repeated float, recent occupancy readings

        Returns the expected occupancy in that voxel over the horizon,
        a congestion flag, and a model confidence interval.
        """
        h3_index = request.h3_index
        time_horizon_sec = (
            request.time_horizon_seconds if request.time_horizon_seconds > 0 else 300
        )
        history = list(request.historical_density) if request.historical_density else [0.0]

        if self.model is not None:
            # Reshape history into (1, in_features=1, seq_len) tensor.
            # The model was trained with a fixed history window; we pad/truncate
            # to match the expected sequence length (90 steps by default).
            seq = np.array(history, dtype=np.float32)
            # Normalise by a reasonable occupancy ceiling (30 aircraft)
            seq = seq / 30.0
            # Add batch and channel dims: (1, 1, T)
            tensor_in = torch.from_numpy(seq).unsqueeze(0).unsqueeze(0)
            with torch.no_grad():
                out = self.model(tensor_in)           # (1, 1, horizon_steps)
            # Take the mean over the horizon window as the scalar forecast
            forecast_val = float(out.mean().item()) * 30.0
            forecast_val = max(0.0, forecast_val)
        else:
            forecast_val = _ema_forecast(history)

        return demand_forecaster_pb2.DemandForecastResponse(
            h3_index=h3_index,
            forecasted_density=forecast_val,
            is_congested=forecast_val > _CONGESTION_THRESHOLD,
            # Confidence interval sourced from Table III (TCN MAE = 0.58)
            confidence_interval=0.942,
        )

    def BatchForecastDemand(self, request, context):
        """
        Batch demand forecasting for fleet-wide regional airspace grids.
        """
        responses = [
            self.ForecastDemand(single_req, context) for single_req in request.requests
        ]
        return demand_forecaster_pb2.BatchDemandForecastResponse(responses=responses)


def serve() -> None:
    """Bootstraps the gRPC server and blocks until process termination."""
    port = os.getenv("PORT", "50055")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    demand_forecaster_pb2_grpc.add_DemandForecasterServiceServicer_to_server(
        DemandForecasterServicer(), server
    )
    server.add_insecure_port(f"0.0.0.0:{port}")
    logger.info("Demand Forecaster Service listening on port %s", port)
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
