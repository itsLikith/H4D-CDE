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
Haversine and spherical geodesic utilities for trajectory calculation.
Computes great-circle distance and intermediate spherical waypoints.
"""

import math

EARTH_RADIUS_KM = 6371.0088


def haversine_km(lat1: float, lon1: float, lat2: float, lon2: float) -> float:
    """
    Computes great-circle distance between two (lat, lon) coordinates in kilometers
    using the Haversine formula.
    """
    p1, p2 = math.radians(lat1), math.radians(lat2)
    dphi = math.radians(lat2 - lat1)
    dlmb = math.radians(lon2 - lon1)

    a = math.sin(dphi / 2.0) ** 2 + math.cos(p1) * math.cos(p2) * math.sin(dlmb / 2.0) ** 2
    c = 2.0 * math.atan2(math.sqrt(a), math.sqrt(1.0 - a))
    return EARTH_RADIUS_KM * c


def km(lat1: float, lon1: float, lat2: float, lon2: float) -> float:
    """Convenience alias for haversine_km."""
    return haversine_km(lat1, lon1, lat2, lon2)


def interpolate(lat1: float, lon1: float, lat2: float, lon2: float, frac: float) -> tuple[float, float]:
    """
    Spherical great-circle interpolation between two points for fraction frac in [0, 1].
    Returns (lat, lon) in degrees.
    """
    if frac <= 0.0:
        return lat1, lon1
    if frac >= 1.0:
        return lat2, lon2

    p1, l1 = math.radians(lat1), math.radians(lon1)
    p2, l2 = math.radians(lat2), math.radians(lon2)

    # Angular distance d between points
    d = 2.0 * math.asin(math.sqrt(
        math.sin((p2 - p1) / 2.0) ** 2 +
        math.cos(p1) * math.cos(p2) * math.sin((l2 - l1) / 2.0) ** 2
    ))

    if d == 0.0:
        return lat1, lon1

    a = math.sin((1.0 - frac) * d) / math.sin(d)
    b = math.sin(frac * d) / math.sin(d)

    x = a * math.cos(p1) * math.cos(l1) + b * math.cos(p2) * math.cos(l2)
    y = a * math.cos(p1) * math.sin(l1) + b * math.cos(p2) * math.sin(l2)
    z = a * math.sin(p1) + b * math.sin(p2)

    lat_rad = math.atan2(z, math.sqrt(x * x + y * y))
    lon_rad = math.atan2(y, x)

    return math.degrees(lat_rad), math.degrees(lon_rad)
