"""
7-feature vector
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
    n_entities,
    closure_rate,
    heading_diff,
    local_density,
    sector_load_forecast,
    wind_shear,
    visibility_km,
) -> list[float]:
    return [
        n_entities,
        closure_rate,
        heading_diff,
        local_density,
        sector_load_forecast,
        wind_shear,
        visibility_km,
    ]
