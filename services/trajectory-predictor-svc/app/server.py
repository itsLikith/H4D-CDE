"""
Serves the trained model over gRPC. This is the service's actual production
interface -- voxel-engine (Go) calls RefineTrajectory exactly as if it were
calling a local function, but it's really a network call to this process.
"""

import logging
from concurrent import futures

import grpc
import joblib
import numpy as np

from gen import trajectory_predictor_pb2 as pb2
from gen import trajectory_predictor_pb2_grpc as pb2_grpc
from gen import common_pb2

# generated from proto/common.proto (Part 8.2) — the same TrajectoryPoint message every language shares
from .features import build_feature_vector, air_density
from . import haversine  # small shared helper, same great-circle math used elsewhere


class TrajectoryPredictorServicer(pb2_grpc.TrajectoryPredictorServiceServicer):
    def __init__(self, model_path: str):
        self._model = joblib.load(model_path)

    def RefineTrajectory(self, request, context):
        fpl = request.flight_plan
        points = []
        waypoints = list(fpl.waypoints)
        for i in range(len(waypoints) - 1):
            a, b = waypoints[i], waypoints[i + 1]
            d_km = haversine.km(a.lat, a.lon, b.lat, b.lon)
            x = np.array(
                [
                    [
                        *build_feature_vector(
                            great_circle_distance_km=d_km,
                            altitude_diff_ft=0.0,
                            wind_speed_kt=request.wind_speed_kt,
                            wind_direction_deg=request.wind_direction_deg,
                            max_accel_mps2=request.max_accel_mps2,
                            cruise_speed_kt=fpl.cruise_speed_kt,
                            altitude_ft=fpl.cruise_altitude_ft,
                        )
                    ]
                ]
            )
            correction_m = float(self._model.predict(x)[0])
            # apply correction_m to the naive interpolated point, append to points...
            points.append(
                common_pb2.TrajectoryPoint(
                    entity_id=fpl.entity_id,
                    lat=a.lat,
                    lon=a.lon,
                    alt_ft=fpl.cruise_altitude_ft,
                )
            )
        return pb2.RefineTrajectoryResponse(points=points)


def serve(model_path: str, port: int = 50051):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=16))
    pb2_grpc.add_TrajectoryPredictorServiceServicer_to_server(
        TrajectoryPredictorServicer(model_path), server
    )
    server.add_insecure_port(
        f"[::]:{port}"
    )  # mTLS added at the mesh/ingress layer, Part 20.6
    server.start()
    logging.info("trajectory-predictor-svc listening on :%d", port)
    server.wait_for_termination()


if __name__ == "__main__":
    serve(model_path="models/trajectory_predictor.joblib")
