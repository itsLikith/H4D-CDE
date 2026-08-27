"""
dilated-causal-convolution TCN as the original design
"""

import torch
import torch.nn as nn


class CausalConv1d(nn.Module):
    def __init__(self, in_ch, out_ch, kernel_size, dilation):
        super().__init__()
        self.left_pad = (kernel_size - 1) * dilation
        self.conv = nn.Conv1d(
            in_ch, out_ch, kernel_size, padding=self.left_pad, dilation=dilation
        )

    def forward(self, x):
        out = self.conv(x)
        return out[:, :, : -self.left_pad] if self.left_pad != 0 else out


class TCNBlock(nn.Module):
    def __init__(self, in_ch, out_ch, kernel_size=3, dilation=1, dropout=0.2):
        super().__init__()
        self.conv1 = CausalConv1d(in_ch, out_ch, kernel_size, dilation)
        self.conv2 = CausalConv1d(out_ch, out_ch, kernel_size, dilation)
        self.relu = nn.ReLU()
        self.dropout = nn.Dropout(dropout)
        self.downsample = nn.Conv1d(in_ch, out_ch, 1) if in_ch != out_ch else None

    def forward(self, x):
        out = self.dropout(self.relu(self.conv1(x)))
        out = self.dropout(self.relu(self.conv2(out)))
        residual = x if self.downsample is None else self.downsample(x)
        return self.relu(out + residual)


class DemandForecasterTCN(nn.Module):
    def __init__(self, num_voxels, hidden_channels=(32, 32, 32), horizon_steps=90):
        super().__init__()
        layers, in_ch = [], num_voxels
        for i, out_ch in enumerate(hidden_channels):
            layers.append(TCNBlock(in_ch, out_ch, dilation=2**i))
            in_ch = out_ch
        self.tcn = nn.Sequential(*layers)
        self.head = nn.Conv1d(in_ch, num_voxels, kernel_size=1)
        self.horizon_steps = horizon_steps

    def forward(self, occupancy_history: torch.Tensor) -> torch.Tensor:
        return self.head(self.tcn(occupancy_history))[:, :, -self.horizon_steps :]
