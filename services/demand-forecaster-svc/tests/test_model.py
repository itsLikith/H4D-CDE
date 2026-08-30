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

import torch
import pytest
from app.model import DemandForecasterTCN


def test_tcn_forward_shape():
    model = DemandForecasterTCN(
        in_features=1, hidden_channels=(16, 16), horizon_steps=90
    )
    x = torch.randn(4, 1, 180)  # batch of 4 sequences of length 180
    out = model(x)
    assert out.shape == (4, 1, 90)
