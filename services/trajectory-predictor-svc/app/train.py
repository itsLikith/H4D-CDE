import joblib
import numpy as np
from sklearn.ensemble import GradientBoostingRegressor
from sklearn.model_selection import train_test_split
from sklearn.metrics import mean_absolute_error

PAPER_BENCHMARK_MAE_M = 15.2


def train_trajectory_predictor(X: np.ndarray, y: np.ndarray, model_out_path: str):
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
    mae_m = mean_absolute_error(y_test, model.predict(X_test))  # Eq. (18)
    print(
        f"Validation MAE: {mae_m:.2f} m  (paper benchmark: {PAPER_BENCHMARK_MAE_M} m)"
    )
    joblib.dump(model, model_out_path)
    return model, mae_m
