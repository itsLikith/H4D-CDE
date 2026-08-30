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
Training pipeline for Risk Scorer (Module 4).
Trains a gradient boosted decision tree classifier with isotonic calibration
to output well-calibrated loss-of-separation risk probabilities in [0.0, 1.0].

Validates Area Under the ROC Curve (AUC-ROC) against the ICSPIS 2025 paper benchmark of 0.89.
"""

import os
import joblib
import numpy as np
from sklearn.ensemble import GradientBoostingClassifier
from sklearn.model_selection import train_test_split
from sklearn.metrics import roc_auc_score
from sklearn.calibration import CalibratedClassifierCV

from .features import build_feature_vector

PAPER_BENCHMARK_AUC = 0.89


def generate_synthetic_conflict_dataset(
    n_samples: int = 6000, random_seed: int = 42
) -> tuple[np.ndarray, np.ndarray]:
    """
    Generates synthetic airspace conflict feature vectors with realistic operational correlation.
    """
    rng = np.random.default_rng(random_seed)

    n_entities = rng.integers(2, 6, size=n_samples)
    closure_rates = rng.uniform(0.0, 80.0, size=n_samples)  # m/s
    heading_diffs = rng.uniform(0.0, 180.0, size=n_samples)  # head-on vs parallel
    local_densities = rng.uniform(1.0, 12.0, size=n_samples)
    sector_loads = rng.uniform(1.0, 25.0, size=n_samples)
    wind_shears = rng.uniform(0.0, 15.0, size=n_samples)
    visibilities = rng.uniform(0.5, 20.0, size=n_samples)

    X_list = []
    y_list = []

    for i in range(n_samples):
        feat = build_feature_vector(
            n_entities_in_conflict=n_entities[i],
            closure_rate_mps=closure_rates[i],
            heading_diff_deg=heading_diffs[i],
            local_traffic_density=local_densities[i],
            sector_load_forecast=sector_loads[i],
            wind_shear_kt_per_100ft=wind_shears[i],
            visibility_km=visibilities[i],
        )
        X_list.append(feat)

        # Logit equation modeling severe loss-of-separation risk
        head_on_factor = np.sin(np.radians(heading_diffs[i] / 2.0))
        risk_logit = (
            0.08 * closure_rates[i]
            + 1.8 * head_on_factor
            + 0.35 * local_densities[i]
            + 0.15 * (n_entities[i] - 2)
            + 0.12 * wind_shears[i]
            - 0.18 * visibilities[i]
            - 4.2
        )
        prob = 1.0 / (1.0 + np.exp(-risk_logit))
        label = int(rng.uniform(0.0, 1.0) < prob)
        y_list.append(label)

    return np.array(X_list, dtype=np.float32), np.array(y_list, dtype=np.int32)


def train_risk_scorer(
    X: np.ndarray,
    y: np.ndarray,
    model_out_path: str = "models/risk_scorer.joblib",
    calibrate: bool = True,
) -> tuple[object, float]:
    """
    Fits a gradient boosted classification model with isotonic calibration on the conflict feature matrix.
    """
    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.2, stratify=y, random_state=42
    )

    base_model = GradientBoostingClassifier(
        n_estimators=300,
        max_depth=4,
        learning_rate=0.05,
        subsample=0.8,
        random_state=42,
    )
    base_model.fit(X_train, y_train)

    if calibrate:
        model = CalibratedClassifierCV(base_model, method="isotonic", cv=3)
        model.fit(X_train, y_train)
    else:
        model = base_model

    y_pred_proba = model.predict_proba(X_test)[:, 1]
    auc = float(roc_auc_score(y_test, y_pred_proba))

    print(
        f"[*] Risk Scorer Validation AUC-ROC: {auc:.3f} (Paper Target: ≥ {PAPER_BENCHMARK_AUC})"
    )

    os.makedirs(os.path.dirname(os.path.abspath(model_out_path)), exist_ok=True)
    joblib.dump(model, model_out_path)
    print(f"[+] Serialized model saved to {model_out_path}")

    return model, auc


if __name__ == "__main__":
    X, y = generate_synthetic_conflict_dataset()
    train_risk_scorer(X, y)
