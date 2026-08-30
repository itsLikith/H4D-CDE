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

import pytest
from app.features import build_feature_vector, air_density
from app.haversine import haversine_km, interpolate


def test_air_density_isa():
    rho_sea_level = air_density(0.0)
    assert pytest.approx(rho_sea_level, rel=1e-2) == 1.225
    rho_5000ft = air_density(5000.0)
    assert rho_5000ft < rho_sea_level


def test_feature_vector_dimension():
    feat = build_feature_vector(
        great_circle_distance_km=2.5,
        altitude_diff_ft=500.0,
        wind_speed_kt=10.0,
        wind_direction_deg=45.0,
        max_accel_mps2=2.5,
        cruise_speed_kt=90.0,
        altitude_ft=1500.0,
    )
    assert len(feat) == 7
    assert isinstance(feat, list)


def test_haversine_distance():
    # DXB to AUH (~115 km)
    dist = haversine_km(25.2532, 55.3657, 24.4330, 54.6511)
    assert 110.0 <= dist <= 125.0


def test_great_circle_interpolate():
    lat, lon = interpolate(25.0, 55.0, 25.0, 56.0, 0.5)
    assert pytest.approx(lat, abs=0.1) == 25.0
    assert pytest.approx(lon, abs=0.1) == 55.5
