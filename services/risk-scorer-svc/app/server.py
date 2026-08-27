import logging
from concurrent import futures

import grpc
import joblib
import numpy as np

from gen import risk_scorer_pb2 as pb2
from gen import risk_scorer_pb2_grpc as pb2_grpc
from .features import build_feature_vector

DEFAULT_ADVISORY_THRESHOLD = 0.5


class RiskScorerServicer(pb2_grpc.RiskScorerServiceServicer):
    def __init__(self, model_path: str):
        self._model = joblib.load(model_path)

    def ScoreConflict(self, request, context):
        x = np.array(
            [
                build_feature_vector(
                    request.n_entities_in_conflict,
                    request.closure_rate_mps,
                    request.heading_diff_deg,
                    request.local_traffic_density,
                    request.sector_load_forecast,
                    request.wind_shear_kt_per_100ft,
                    request.visibility_km,
                )
            ]
        )
        score = float(self._model.predict_proba(x)[0][1])
        return pb2.ScoreConflictResponse(risk_score=score)


def serve(model_path: str, port: int = 50052):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=32))
    # sized higher than
    # trajectory-predictor-svc because this service sits on every
    # single conflict, not just once per flight-plan submission
    pb2_grpc.add_RiskScorerServiceServicer_to_server(
        RiskScorerServicer(model_path), server
    )
    server.add_insecure_port(f"[::]:{port}")
    server.start()
    logging.info("risk-scorer-svc listening on :%d", port)
    server.wait_for_termination()
