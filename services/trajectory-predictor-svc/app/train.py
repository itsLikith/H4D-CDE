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
Training pipeline for the Trajectory Predictor (Module 2).
Trains a Gradient Boosting Regressor (GBM) to predict physically accurate
trajectory cross-track and along-track deviations from nominal flight plans.

Evaluates Mean Absolute Error (MAE, Eq. 18):
    MAE = (1/n) * Σ |y_i - ŷ_i|
Validates against the ICSPIS 2025 paper benchmark of MAE = 15.2 m.
"""

import os
import joblib
import numpy as np
from sklearn.ensemble import GradientBoostingRegressor
from sklearn.model_selection import train_test_split
from sklearn.metrics import mean_absolute_error

from .features import build_feature_vector

PAPER_BENCHMARK_MAE_M = 15.2


def generate_synthetic_training_data(
    n_samples: int = 5000, random_seed: int = 42
) -> tuple[np.ndarray, np.ndarray]:
    """
    Generates physically consistent flight training vectors under varying
    atmospheric and operational conditions.
    """
    rng = np.random.default_rng(random_seed)

    distances_km = rng.uniform(5.0, 150.0, size=n_samples)
    alt_diffs_ft = rng.uniform(-1000.0, 1000.0, size=n_samples)
    wind_speeds_kt = rng.uniform(0.0, 45.0, size=n_samples)
    wind_dirs_deg = rng.uniform(0.0, 360.0, size=n_samples)
    max_accels_mps2 = rng.uniform(1.0, 4.0, size=n_samples)
    cruise_speeds_kt = rng.uniform(50.0, 130.0, size=n_samples)
    altitudes_ft = rng.uniform(300.0, 3000.0, size=n_samples)

    X_list = []
    y_list = []

    for i in range(n_samples):
        feat = build_feature_vector(
            great_circle_distance_km=distances_km[i],
            altitude_diff_ft=alt_diffs_ft[i],
            wind_speed_kt=wind_speeds_kt[i],
            wind_direction_deg=wind_dirs_deg[i],
            max_accel_mps2=max_accels_mps2[i],
            cruise_speed_kt=cruise_speeds_kt[i],
            altitude_ft=altitudes_ft[i],
        )
        X_list.append(feat)

        # Physics-informed ground truth deviation:
        # Crosswind drift component + density correction + acceleration lag
        headwind_comp = wind_speeds_kt[i] * np.cos(np.radians(wind_dirs_deg[i]))
        crosswind_comp = wind_speeds_kt[i] * np.sin(np.radians(wind_dirs_deg[i]))
        true_dev_m = (
            0.35 * crosswind_comp
            + 0.08 * headwind_comp
            + 0.02 * (alt_diffs_ft[i] / 100.0)
            + 0.5 * (1.225 - feat[-1]) * 10.0
            + rng.normal(0.0, 12.0)
        )
        y_list.append(true_dev_m)

    return np.array(X_list, dtype=np.float32), np.array(y_list, dtype=np.float32)


def train_trajectory_predictor(
    X: np.ndarray,
    y: np.ndarray,
    model_out_path: str = "models/trajectory_predictor.joblib",
) -> tuple[GradientBoostingRegressor, float]:
    """
    Fits a Gradient Boosting Regressor (GBM) on the 7-feature dataset.
    Saves the serialized model artifact using joblib.
    """
    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.2, random_state=42
    )

    model = GradientBoostingRegressor(
        n_estimators=300,
        max_depth=4,
        learning_rate=0.05,
        subsample=0.8,
        random_state=42,
    )
    model.fit(X_train, y_train)

    y_pred = model.predict(X_test)
    mae_m = float(mean_absolute_error(y_test, y_pred))

    print(
        f"[*] Trajectory Predictor Validation MAE: {mae_m:.2f} m (Paper Target: ≤ {PAPER_BENCHMARK_MAE_M} m)"
    )

    os.makedirs(os.path.dirname(os.path.abspath(model_out_path)), exist_ok=True)
    joblib.dump(model, model_out_path)
    print(f"[+] Serialized model saved to {model_out_path}")

    return model, mae_m


if __name__ == "__main__":
    X, y = generate_synthetic_training_data()
    train_trajectory_predictor(X, y)
