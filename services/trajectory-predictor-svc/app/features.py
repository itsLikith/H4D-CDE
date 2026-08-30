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
Feature Engineering for Trajectory Predictor (Module 2).
Implements the 7-feature vector from Equation (8) of Sahadevan et al. (ICSPIS 2025):

    ŷ = f_GBM(X)
    X = [d, Δh, ws, wd, a_max, v_cruise, ρ_air]

Where:
    - d: Great-circle distance between consecutive waypoints (km)
    - Δh: Altitude differential (ft)
    - ws: Wind speed (knots)
    - wd: Wind direction (degrees 0-360)
    - a_max: Maximum vehicle acceleration (m/s^2)
    - v_cruise: Cruise speed (knots)
    - ρ_air: Air density computed from barometric standard atmosphere formula (kg/m^3)
"""

SEA_LEVEL_DENSITY_KGM3 = 1.225
TEMP_LAPSE_RATE_K_PER_M = 0.0065
SEA_LEVEL_TEMP_K = 288.15


def air_density(alt_ft: float) -> float:
    """
    Computes ambient air density ρ_air at a given altitude using the
    International Standard Atmosphere (ISA) barometric lapse formula:
        T(h) = T0 - L * h
        ρ(h) = ρ0 * (T(h) / T0)^(g / (R*L) - 1) ≈ ρ0 * (T / T0)^4.2559
    """
    alt_m = max(alt_ft, 0.0) * 0.3048
    temp_k = max(SEA_LEVEL_TEMP_K - TEMP_LAPSE_RATE_K_PER_M * alt_m, 180.0)
    return SEA_LEVEL_DENSITY_KGM3 * (temp_k / SEA_LEVEL_TEMP_K) ** 4.2559


def build_feature_vector(
    great_circle_distance_km: float,
    altitude_diff_ft: float,
    wind_speed_kt: float,
    wind_direction_deg: float,
    max_accel_mps2: float,
    cruise_speed_kt: float,
    altitude_ft: float,
) -> list[float]:
    """
    Constructs the 7-dimensional input feature vector for the Gradient Boosting Regressor (GBM).
    """
    return [
        float(great_circle_distance_km),
        float(altitude_diff_ft),
        float(wind_speed_kt),
        float(wind_direction_deg),
        float(max_accel_mps2),
        float(cruise_speed_kt),
        float(air_density(altitude_ft)),
    ]
