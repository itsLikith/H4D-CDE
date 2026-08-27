"""
Runs continuously: every FORECAST_INTERVAL_S, pull recent occupancy history
from Redis (written by voxel-engine), run one forward pass, and
write the forecast back to Redis under a key the Adaptive Discretization
Engine reads directly -- no RPC needed for this path.
"""

import time
import json
import redis
import torch

from .model import DemandForecasterTCN

FORECAST_INTERVAL_S = 30
HORIZON_STEPS = 90  # 15 minutes at 10s bins


def run(model: DemandForecasterTCN, redis_client: redis.Redis):
    while True:
        active_cells = redis_client.smembers(
            "active_voxel_cells"
        )  # maintained by voxel-engine
        for cell in active_cells:
            history = _load_occupancy_history(redis_client, cell)
            with torch.no_grad():
                forecast = model(history).squeeze().tolist()
            redis_client.set(
                f"forecast:{cell.decode()}",
                json.dumps(forecast),
                ex=FORECAST_INTERVAL_S * 3,
            )
        time.sleep(FORECAST_INTERVAL_S)


def _load_occupancy_history(
    redis_client: redis.Redis, cell: bytes
) -> torch.Tensor: ...  # pull the last k occupancy counts for this cell from a Redis time series / sorted set
