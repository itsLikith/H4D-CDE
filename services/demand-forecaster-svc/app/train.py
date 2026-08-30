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
Training pipeline for Demand Forecaster (Module 3).
Trains the Dilated Causal TCN on voxel occupancy time series.
Evaluates Mean Absolute Percentage Error (MAPE, Eq. 19):

    MAPE = (1/n) * Σ |(O_i - Ô_i) / O_i| * 100%

Validates against the ICSPIS 2025 paper benchmark of MAPE = 8.70%.
"""

import os
import torch
import torch.nn as nn
import torch.optim as optim
import numpy as np

from .model import DemandForecasterTCN

PAPER_BENCHMARK_MAPE_PCT = 8.70


def generate_synthetic_occupancy_series(
    n_series: int = 100,
    total_steps: int = 300,
    random_seed: int = 42,
) -> np.ndarray:
    """
    Generates synthetic airspace voxel occupancy time series exhibiting diurnal waves,
    bursty peak demand spikes, and smooth trends.
    """
    rng = np.random.default_rng(random_seed)
    t = np.linspace(0, 6 * np.pi, total_steps)

    series_list = []
    for _ in range(n_series):
        phase = rng.uniform(0, 2 * np.pi)
        amplitude = rng.uniform(4.0, 10.0)
        base_demand = 8.0 + amplitude * np.sin(t + phase) + 2.0 * np.cos(2 * t + phase)
        bursts = rng.poisson(lam=0.1, size=total_steps) * rng.uniform(1.0, 2.0, size=total_steps)
        noise = rng.normal(0.0, 0.15, size=total_steps)
        occ = np.maximum(base_demand + bursts + noise, 2.0)
        series_list.append(occ)

    return np.array(series_list, dtype=np.float32)


def train_demand_forecaster(
    data: np.ndarray,
    history_len: int = 180,
    horizon_steps: int = 90,
    epochs: int = 40,
    model_out_path: str = "models/demand_forecaster.pt",
) -> tuple[DemandForecasterTCN, float]:
    """
    Trains the TCN model and saves PyTorch model weights.
    """
    # Create sliding window dataset
    X_list, y_list = [], []
    for s in data:
        for start in range(0, len(s) - history_len - horizon_steps + 1, 10):
            X_list.append(s[start : start + history_len])
            y_list.append(s[start + history_len : start + history_len + horizon_steps])

    X_t = torch.tensor(np.array(X_list), dtype=torch.float32).unsqueeze(1)
    y_t = torch.tensor(np.array(y_list), dtype=torch.float32).unsqueeze(1)

    split_idx = int(0.85 * len(X_t))
    X_train, X_val = X_t[:split_idx], X_t[split_idx:]
    y_train, y_val = y_t[:split_idx], y_t[split_idx:]

    model = DemandForecasterTCN(in_features=1, hidden_channels=(32, 64, 32), horizon_steps=horizon_steps, dropout=0.05)
    criterion = nn.L1Loss()
    optimizer = optim.Adam(model.parameters(), lr=0.003, weight_decay=1e-6)
    scheduler = optim.lr_scheduler.CosineAnnealingLR(optimizer, T_max=epochs)

    dataset = torch.utils.data.TensorDataset(X_train, y_train)
    loader = torch.utils.data.DataLoader(dataset, batch_size=32, shuffle=True)

    model.train()
    for _ in range(epochs):
        for bx, by in loader:
            optimizer.zero_grad()
            pred = model(bx)
            loss = criterion(pred, by)
            loss.backward()
            optimizer.step()
        scheduler.step()

    # Evaluate validation MAPE (Eq. 19)
    model.eval()
    with torch.no_grad():
        val_pred = model(X_val)
        actual = y_val.numpy()
        predicted = np.maximum(val_pred.numpy(), 0.0)
        mape_pct = float(np.mean(np.abs((actual - predicted) / np.maximum(actual, 1.0))) * 100.0)

    print(f"[*] Demand Forecaster Validation MAPE: {mape_pct:.2f}% (Paper Target: ≤ {PAPER_BENCHMARK_MAPE_PCT}%)")

    os.makedirs(os.path.dirname(os.path.abspath(model_out_path)), exist_ok=True)
    torch.save(model.state_dict(), model_out_path)
    print(f"[+] Serialized model weights saved to {model_out_path}")

    return model, mape_pct


if __name__ == "__main__":
    raw_data = generate_synthetic_occupancy_series()
    train_demand_forecaster(raw_data)
