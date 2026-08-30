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
Temporal Convolutional Network (TCN) for Demand Forecasting (Module 3).
Implements Equation (9) of Sahadevan et al. (ICSPIS 2025):

    Ô(t + Δt) = f_TCN(O(t - k : t), θ)

Utilizes Dilated Causal 1D Convolutions with residual connections (Bai et al., 2018)
to forecast future voxel occupancy over a 15-minute horizon (90 steps of 10s bins)
without temporal leakage.
"""

import torch
import torch.nn as nn


class CausalConv1d(nn.Module):
    """
    1D Causal Convolution with Left-Padding.
    Ensures that prediction at step t depends only on inputs at steps <= t (no lookahead bias).
    """

    def __init__(self, in_channels: int, out_channels: int, kernel_size: int, dilation: int):
        super().__init__()
        self.padding = (kernel_size - 1) * dilation
        self.conv = nn.Conv1d(
            in_channels,
            out_channels,
            kernel_size=kernel_size,
            padding=self.padding,
            dilation=dilation,
        )

    def forward(self, x: torch.Tensor) -> torch.Tensor:
        # x shape: (batch_size, in_channels, seq_len)
        out = self.conv(x)
        if self.padding != 0:
            return out[:, :, : -self.padding]
        return out


class TCNBlock(nn.Module):
    """
    Residual TCN Block with 2 Dilated Causal Convolutions, ReLU, Dropout, and optional 1x1 projection.
    """

    def __init__(
        self,
        in_channels: int,
        out_channels: int,
        kernel_size: int = 3,
        dilation: int = 1,
        dropout: float = 0.2,
    ):
        super().__init__()
        self.conv1 = CausalConv1d(in_channels, out_channels, kernel_size, dilation)
        self.relu1 = nn.ReLU()
        self.dropout1 = nn.Dropout(dropout)

        self.conv2 = CausalConv1d(out_channels, out_channels, kernel_size, dilation)
        self.relu2 = nn.ReLU()
        self.dropout2 = nn.Dropout(dropout)

        self.downsample = (
            nn.Conv1d(in_channels, out_channels, kernel_size=1)
            if in_channels != out_channels
            else None
        )
        self.relu_out = nn.ReLU()

    def forward(self, x: torch.Tensor) -> torch.Tensor:
        residual = x if self.downsample is None else self.downsample(x)
        out = self.dropout1(self.relu1(self.conv1(x)))
        out = self.dropout2(self.relu2(self.conv2(out)))
        return self.relu_out(out + residual)


class DemandForecasterTCN(nn.Module):
    """
    Multi-scale Dilated Temporal Convolutional Network.
    Takes a time-series history of voxel occupancy and outputs multi-step ahead forecasts.
    """

    def __init__(
        self,
        in_features: int = 1,
        hidden_channels: tuple[int, ...] = (32, 32, 32),
        horizon_steps: int = 90,
        kernel_size: int = 3,
        dropout: float = 0.15,
    ):
        super().__init__()
        self.horizon_steps = horizon_steps

        layers = []
        in_ch = in_features
        for i, out_ch in enumerate(hidden_channels):
            dilation = 2 ** i
            layers.append(
                TCNBlock(
                    in_channels=in_ch,
                    out_channels=out_ch,
                    kernel_size=kernel_size,
                    dilation=dilation,
                    dropout=dropout,
                )
            )
            in_ch = out_ch

        self.network = nn.Sequential(*layers)
        self.head = nn.Conv1d(in_ch, in_features, kernel_size=1)

    def forward(self, x: torch.Tensor) -> torch.Tensor:
        # Expected input shape: (batch, in_features, history_len)
        features = self.network(x)
        out = self.head(features)
        # Return the last `horizon_steps` sequence
        return out[:, :, -self.horizon_steps :]
