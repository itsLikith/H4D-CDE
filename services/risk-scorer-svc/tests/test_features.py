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
from app.features import build_feature_vector


def test_risk_scorer_feature_vector_dimension():
    feat = build_feature_vector(
        n_entities_in_conflict=2,
        closure_rate_mps=45.0,
        heading_diff_deg=180.0,
        local_traffic_density=4.0,
        sector_load_forecast=8.0,
        wind_shear_kt_per_100ft=3.5,
        visibility_km=10.0,
    )
    assert len(feat) == 7
    assert isinstance(feat, list)
