"""Builds the 7-feature vector from Eq. (8) — same math as the original
single-service design, now living inside its own service package."""

import math

SEA_LEVEL_DENSITY_KGM3 = 1.225
TEMP_LAPSE_RATE_K_PER_M = 0.0065
SEA_LEVEL_TEMP_K = 288.15


def air_density(alt_ft: float) -> float:
    alt_m = alt_ft * 0.3048
    temp_k = SEA_LEVEL_TEMP_K - TEMP_LAPSE_RATE_K_PER_M * alt_m
    return SEA_LEVEL_DENSITY_KGM3 * (temp_k / SEA_LEVEL_TEMP_K) ** 4.2559


def build_feature_vector(
    great_circle_distance_km,
    altitude_diff_ft,
    wind_speed_kt,
    wind_direction_deg,
    max_accel_mps2,
    cruise_speed_kt,
    altitude_ft,
) -> list[float]:
    return [
        great_circle_distance_km,
        altitude_diff_ft,
        wind_speed_kt,
        wind_direction_deg,
        max_accel_mps2,
        cruise_speed_kt,
        air_density(altitude_ft),
    ]
