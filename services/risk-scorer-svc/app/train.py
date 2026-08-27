import joblib
import xgboost as xgb
from sklearn.model_selection import train_test_split
from sklearn.metrics import roc_auc_score
from sklearn.calibration import CalibratedClassifierCV

PAPER_BENCHMARK_AUC = 0.89


def train_risk_scorer(X, y, model_out_path: str, calibrate: bool = True):
    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.2, stratify=y, random_state=42
    )
    base_model = xgb.XGBClassifier(
        n_estimators=400,
        max_depth=5,
        learning_rate=0.05,
        subsample=0.8,
        colsample_bytree=0.8,
        eval_metric="auc",
        random_state=42,
    )
    base_model.fit(X_train, y_train)
    model = base_model
    if calibrate:
        model = CalibratedClassifierCV(base_model, method="isotonic", cv="prefit")
        model.fit(X_test, y_test)
    auc = roc_auc_score(y_test, model.predict_proba(X_test)[:, 1])
    print(f"Validation AUC-ROC: {auc:.2f}  (paper benchmark: {PAPER_BENCHMARK_AUC})")
    joblib.dump(model, model_out_path)
    return model, auc
