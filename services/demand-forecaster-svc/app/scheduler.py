# Copyright 2026 Likith Saragadam
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
Background Forecasting Daemon for Demand Forecaster (Module 3).
Continuously evaluates near-term voxel demand and writes predictions directly into Redis
under keys 'forecast:<h3_cell>' for direct, zero-RPC reads by VoxelEngine's Adaptive Discretization Engine.
"""

import os
import time
import json
import logging
import redis
import torch

from .model import DemandForecasterTCN

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")

FORECAST_INTERVAL_S = int(os.getenv("FORECAST_INTERVAL_S", "30"))
HORIZON_STEPS = int(os.getenv("FORECAST_HORIZON_STEPS", "90"))  # 15 min at 10s bins
HISTORY_LEN = 180


def load_occupancy_history(redis_client: redis.Redis, cell: str) -> torch.Tensor:
    """
    Pulls recent occupancy history from Redis sorted set or list.
    Falls back to baseline traffic profile if history is sparse.
    """
    try:
        raw_vals = redis_client.lrange(f"history:{cell}", -HISTORY_LEN, -1)
        if raw_vals and len(raw_vals) >= 10:
            floats = [float(v.decode()) for v in raw_vals]
            if len(floats) < HISTORY_LEN:
                floats = [floats[0]] * (HISTORY_LEN - len(floats)) + floats
            return torch.tensor(floats, dtype=torch.float32).view(1, 1, HISTORY_LEN)
    except Exception as e:
        logging.debug("Could not read Redis history for %s: %s", cell, e)

    # Default synthetic historical baseline for initial cold start
    base = [2.0] * HISTORY_LEN
    return torch.tensor(base, dtype=torch.float32).view(1, 1, HISTORY_LEN)


def run_scheduler(model: DemandForecasterTCN, redis_client: redis.Redis):
    """
    Continuous background loop polling active cells and refreshing forecast keys in Redis.
    """
    logging.info("[*] Starting demand forecasting loop (interval: %d s)...", FORECAST_INTERVAL_S)
    model.eval()

    while True:
        try:
            active_cells = redis_client.smembers("active_voxel_cells")
            if not active_cells:
                # Poll default reference airspace cells if none active
                active_cells = {b"8828308281fffff", b"8828308285fffff", b"8828308287fffff"}

            for cell_bytes in active_cells:
                cell_str = cell_bytes.decode() if isinstance(cell_bytes, bytes) else str(cell_bytes)
                history = load_occupancy_history(redis_client, cell_str)

                with torch.no_grad():
                    pred = model(history).squeeze().tolist()
                    if isinstance(pred, float):
                        pred = [int(round(max(pred, 0.0)))]
                    else:
                        pred = [int(round(max(p, 0.0))) for p in pred]

                redis_key = f"forecast:{cell_str}"
                redis_client.set(redis_key, json.dumps(pred), ex=FORECAST_INTERVAL_S * 4)

            logging.debug("Refreshed demand forecasts for %d cells", len(active_cells))
        except Exception as e:
            logging.warning("Scheduler cycle encountered error: %s", e)

        time.sleep(FORECAST_INTERVAL_S)
