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
Risk Scorer gRPC Service (Module 4).

Implements dynamic collision risk scoring as described in Sahadevan et al.
(ICSPIS 2025), Section IV-D. An XGBoost classifier outputs a probability
P(conflict) from the 7-dimensional feature vector:

    X' = [n, closure_rate, heading_diff, local_density,
          sector_load_forecast, wind_shear, visibility]

When the pre-trained model is unavailable the service falls back to a
Fermi-approximation sigmoid formula derived from the same feature space.
"""

import logging
import math
import os
import sys
from concurrent import futures

import grpc
import joblib
import numpy as np

# ---------------------------------------------------------------------------
# Add the generated protobuf stubs (gen/) to sys.path so they are importable
# inside the Docker container where PYTHONPATH may not be pre-configured.
# ---------------------------------------------------------------------------
_THIS_DIR = os.path.dirname(os.path.abspath(__file__))
_SVC_DIR = os.path.dirname(_THIS_DIR)
_GEN_DIR = os.path.join(_SVC_DIR, "gen")
if _GEN_DIR not in sys.path:
    sys.path.insert(0, _GEN_DIR)

import risk_scorer_pb2             # noqa: E402
import risk_scorer_pb2_grpc        # noqa: E402
from app.features import build_feature_vector  # noqa: E402

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)
logger = logging.getLogger("risk-scorer-svc")

# Risk severity thresholds from Table II of the paper
_CRITICAL_THRESHOLD = 0.80
_MEDIUM_THRESHOLD = 0.50


def _fermi_sigmoid(n: float, closure_rate: float, time_diff: float) -> float:
    """
    Physics-based fallback risk estimate when no trained model is present.
    Uses a logistic function of closure rate and conflict count analogous
    to the Fermi approximation described in Sahadevan et al. Section III-C.
    """
    return 1.0 / (1.0 + math.exp(-0.08 * (closure_rate - 5.0) - 0.3 * (n - 1.0)))


class RiskScorerServicer(risk_scorer_pb2_grpc.RiskScorerServiceServicer):
    """
    gRPC Servicer implementing the RiskScorerService contract.

    Loads an XGBoost model on startup (MODEL_PATH env var).
    Falls back gracefully if the model file is not yet present
    (e.g. before `make train-models` has been run).
    """

    def __init__(self, model_path: str | None = None):
        if model_path is None:
            model_path = os.getenv("MODEL_PATH", "models/risk_xgb.joblib")

        self.model = None
        if os.path.exists(model_path):
            try:
                self.model = joblib.load(model_path)
                logger.info("Loaded XGBoost risk model from %s", model_path)
            except Exception as exc:
                logger.warning("Failed to load model from %s: %s", model_path, exc)
        else:
            logger.warning(
                "Model file not found at %s – using Fermi-sigmoid fallback.", model_path
            )

    # ------------------------------------------------------------------
    # RPC handlers
    # ------------------------------------------------------------------

    def ScoreRisk(self, request, context):
        """
        Scores pairwise collision probability between two 4D trajectories.

        The feature vector is built from request fields that map directly
        onto the 7 features defined in features.build_feature_vector().
        Missing fields default to safe conservative values so the service
        is always callable without crashing.
        """
        feat = build_feature_vector(
            n_entities_in_conflict=getattr(request, "n_entities", 2),
            closure_rate_mps=getattr(request, "closure_rate_mps", 5.0),
            heading_diff_deg=getattr(request, "heading_diff_deg", 45.0),
            local_traffic_density=getattr(request, "local_traffic_density", 3.0),
            sector_load_forecast=getattr(request, "sector_load_forecast", 5.0),
            wind_shear_kt_per_100ft=getattr(request, "wind_shear", 0.0),
            visibility_km=getattr(request, "visibility_km", 10.0),
        )

        if self.model is not None:
            try:
                raw_score = float(self.model.predict_proba([feat])[0][1])
            except Exception as exc:
                logger.warning("Model inference failed: %s – using Fermi fallback.", exc)
                raw_score = _fermi_sigmoid(feat[0], feat[1], 0.0)
        else:
            raw_score = _fermi_sigmoid(feat[0], feat[1], 0.0)

        score = float(np.clip(raw_score, 0.0, 1.0))

        if score > _CRITICAL_THRESHOLD:
            severity = "CRITICAL"
        elif score > _MEDIUM_THRESHOLD:
            severity = "MEDIUM"
        else:
            severity = "LOW"

        return risk_scorer_pb2.ScoreRiskResponse(
            flight_id_a=request.flight_id_a,
            flight_id_b=request.flight_id_b,
            risk_score=score,
            severity=severity,
        )

    def BatchScoreRisk(self, request, context):
        """
        Batch scoring for fleet-wide parallel conflict cluster analysis.
        """
        responses = [
            self.ScoreRisk(pair_req, context) for pair_req in request.requests
        ]
        return risk_scorer_pb2.BatchScoreRiskResponse(responses=responses)


def serve() -> None:
    """Bootstraps the gRPC server and blocks until process termination."""
    port = os.getenv("PORT", "50053")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    risk_scorer_pb2_grpc.add_RiskScorerServiceServicer_to_server(
        RiskScorerServicer(), server
    )
    server.add_insecure_port(f"0.0.0.0:{port}")
    logger.info("Risk Scorer Service listening on port %s", port)
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
