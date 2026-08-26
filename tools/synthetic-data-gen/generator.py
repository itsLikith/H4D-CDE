# Copyright 2026 Likith Saragadam

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# ```
# http://www.apache.org/licenses/LICENSE-2.0
# ```

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


"""
Physics-based synthetic trajectory generator,  with realistic UAS dynamics
"""

import math
from dataclasses import dataclass

EARTH_RADIUS_KM = 6371.0088

# Approximate reference coordinates.
# A real deployment sources exact coordinates from the relevant AIP (Aeronautical Information Publication), not this file.
AIRPORTS = {
    "OMDB": {"name": "Dubai Intl (DXB)",    "lat": 25.2532, "lon": 55.3657},
    "OMAA": {"name": "Abu Dhabi Intl (AUH)", "lat": 24.4330, "lon": 54.6511},
    "OMDW": {"name": "Al Maktoum Intl / DWC", "lat": 24.8964, "lon": 55.1613},
    "OMSJ": {"name": "Sharjah Intl (SHJ)",   "lat": 25.3286, "lon": 55.5172},
}


def _haversine_km(lat1, lon1, lat2, lon2) -> float:
    p1, p2 = math.radians(lat1), math.radians(lat2)
    dphi, dlmb = math.radians(lat2 - lat1), math.radians(lon2 - lon1)
    a = math.sin(dphi / 2) ** 2 + math.cos(p1) * math.cos(p2) * math.sin(dlmb / 2) ** 2
    return 2 * EARTH_RADIUS_KM * math.asin(math.sqrt(a))


def _great_circle_interpolate(lat1, lon1, lat2, lon2, f: float) -> tuple[float, float]:
    p1, l1, p2, l2 = map(math.radians, (lat1, lon1, lat2, lon2))
    d = 2 * math.asin(math.sqrt(
        math.sin((p2 - p1) / 2) ** 2 + math.cos(p1) * math.cos(p2) * math.sin((l2 - l1) / 2) ** 2
    ))
    if d == 0:
        return lat1, lon1
    a, b = math.sin((1 - f) * d) / math.sin(d), math.sin(f * d) / math.sin(d)
    x = a * math.cos(p1) * math.cos(l1) + b * math.cos(p2) * math.cos(l2)
    y = a * math.cos(p1) * math.sin(l1) + b * math.cos(p2) * math.sin(l2)
    z = a * math.sin(p1) + b * math.sin(p2)
    return math.degrees(math.atan2(z, math.sqrt(x * x + y * y))), math.degrees(math.atan2(y, x))


@dataclass
class TrajectoryConfig:
    entity_id: str
    origin_icao: str
    destination_icao: str
    eobt_s: float
    cruise_altitude_ft: float = 1500.0
    cruise_speed_kt: float = 90.0
    climb_rate_fpm: float = 500.0
    sample_dt_s: float = 1.0


def generate_trajectory(cfg: TrajectoryConfig) -> list[dict]:
    origin, dest = AIRPORTS[cfg.origin_icao], AIRPORTS[cfg.destination_icao]
    distance_km = _haversine_km(origin["lat"], origin["lon"], dest["lat"], dest["lon"])
    cruise_speed_kmh = cfg.cruise_speed_kt * 1.852
    cruise_time_s = (distance_km / cruise_speed_kmh) * 3600
    climb_time_s = descent_time_s = (cfg.cruise_altitude_ft / cfg.climb_rate_fpm) * 60
    total_time_s = climb_time_s + cruise_time_s + descent_time_s

    points, t = [], 0.0
    while t <= total_time_s:
        if t < climb_time_s:
            alt_ft, frac = (t / climb_time_s) * cfg.cruise_altitude_ft, 0.0
        elif t > total_time_s - descent_time_s:
            remaining = total_time_s - t
            alt_ft, frac = (remaining / descent_time_s) * cfg.cruise_altitude_ft, 1.0
        else:
            alt_ft, frac = cfg.cruise_altitude_ft, (t - climb_time_s) / cruise_time_s

        lat, lon = _great_circle_interpolate(origin["lat"], origin["lon"], dest["lat"], dest["lon"], frac)
        points.append({"entity_id": cfg.entity_id, "t_s": cfg.eobt_s + t, "lat": lat, "lon": lon, "alt_ft": max(alt_ft, 0.0)})
        t += cfg.sample_dt_s
    return points


def build_paper_reference_scenario() -> list[dict]:
    """3 UAVs across the 4 UAE airports, staggered so paths cross near shared
    airspace -- your first integration-test target (Part 18.2)."""
    configs = [
        TrajectoryConfig("UAV-1", "OMDB", "OMAA", eobt_s=0),
        TrajectoryConfig("UAV-2", "OMSJ", "OMDW", eobt_s=60),
        TrajectoryConfig("UAV-3", "OMDW", "OMSJ", eobt_s=120),
    ]
    return [pt for cfg in configs for pt in generate_trajectory(cfg)]