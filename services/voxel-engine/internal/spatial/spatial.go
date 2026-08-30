// Copyright 2026 Likith Saragadam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package spatial provides geospatial mapping functions based on Uber's H3 Discrete Global Grid System.
package spatial

import (
	"fmt"

	"github.com/uber/h3-go/v4"
)

// DefaultResolution is set to H3 Resolution 8 based on the ICSPIS 2025 paper.
// Resolution 8 provides average hexagon area ≈ 0.737 km² (edge length ≈ 461 m),
// ideal for urban low-altitude airspace indexing.
const DefaultResolution = 8

// PointToH3Cell implements Equation (1) and (3) from Sahadevan et al. (ICSPIS 2025):
//
//	H3_cell = h3.latlng_to_cell(lat, lon, res)
//
// Converts geographic coordinates (latitude, longitude in degrees) into an H3 hexagonal spatial cell index.
func PointToH3Cell(lat, lon float64, resolution int) (h3.Cell, error) {
	if resolution < 0 || resolution > 15 {
		return 0, fmt.Errorf("invalid H3 resolution %d (must be between 0 and 15)", resolution)
	}
	latLng := h3.NewLatLng(lat, lon)
	return h3.LatLngToCell(latLng, resolution)
}

// NeighborCells returns the k-ring grid disk: the origin cell plus all hexagonal neighbors
// within grid distance k.
// For k=1 on a standard hexagonal cell, this returns exactly 7 cells (the central cell + 6 immediate neighbors).
func NeighborCells(cell h3.Cell, k int) ([]h3.Cell, error) {
	if k < 0 {
		return nil, fmt.Errorf("k-ring distance cannot be negative: %d", k)
	}
	return h3.GridDisk(cell, k)
}
