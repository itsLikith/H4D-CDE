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
Feature Engineering for Risk Scorer (Module 4).
Implements the 7-feature vector X' from Equation (10) of Sahadevan et al. (ICSPIS 2025):

    risk = f_XGBoost(X')
    X' = [n, closure_rate, heading_diff, local_density, sector_load_forecast, wind_shear, visibility]

Where:
    - n: Number of entities in conflict volume
    - closure_rate: Relative closure speed between aircraft pairs (m/s)
    - heading_diff: Heading angular divergence (degrees in [0, 180])
    - local_density: Current voxel traffic occupancy count
    - sector_load_forecast: Anticipated downstream sector load from Demand Forecaster (Module 3)
    - wind_shear: Vertical wind speed gradient (knots per 100 ft)
    - visibility: Meteorological visual range (km)
"""

FEATURE_NAMES = [
    "n_entities_in_conflict",
    "closure_rate_mps",
    "heading_diff_deg",
    "local_traffic_density",
    "sector_load_forecast",
    "wind_shear_kt_per_100ft",
    "visibility_km",
]


def build_feature_vector(
    n_entities_in_conflict: int,
    closure_rate_mps: float,
    heading_diff_deg: float,
    local_traffic_density: float,
    sector_load_forecast: float,
    wind_shear_kt_per_100ft: float,
    visibility_km: float,
) -> list[float]:
    """
    Builds the normalized 7-dimensional conflict risk feature vector X'.
    """
    return [
        float(n_entities_in_conflict),
        float(closure_rate_mps),
        float(heading_diff_deg),
        float(local_traffic_density),
        float(sector_load_forecast),
        float(wind_shear_kt_per_100ft),
        float(visibility_km),
    ]
