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

// Package temporal defines altitude and time binning operations and the composite 4D VoxelKey type.
package temporal

import (
	"fmt"
	"math"

	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/spatial"
	"github.com/uber/h3-go/v4"
)

const (
	// AltitudeBinFt specifies the vertical bin resolution (100 ft vertical bands) as per Eq. (4).
	AltitudeBinFt = 100

	// TimeBinS specifies the temporal bin width (10-second sliding windows) as per Eq. (5).
	TimeBinS = 10
)

// VoxelKey represents the composite 4D space-time voxel bucket from Equation (6):
//
//	voxel_key = (h_voxel, a_bin, t_bin)
//
// Where:
//   - H3Cell: Hexagonal spatial column index (Resolution 8)
//   - AltBinFt: Truncated altitude bucket floor in feet
//   - TimeBinS: Truncated time bucket floor in seconds since epoch
type VoxelKey struct {
	H3Cell   h3.Cell
	AltBinFt int
	TimeBinS int
}

// String serializes the composite key into a flat string identifier (e.g. "8828308281fffff:400:36000")
// used as the Redis Set key and protobuf wire representation.
func (k VoxelKey) String() string {
	return fmt.Sprintf("%s:%d:%d", k.H3Cell.String(), k.AltBinFt, k.TimeBinS)
}

// AltitudeBin implements Equation (4):
//
//	ab_i = (alt_i // 100) * 100
//
// Discretizes continuous altitude (ft) into discrete vertical bands.
func AltitudeBin(altFt float64, binWidthFt int) int {
	if binWidthFt <= 0 {
		binWidthFt = AltitudeBinFt
	}
	// Floor to handle non-negative altitudes accurately
	return int(math.Floor(altFt/float64(binWidthFt))) * binWidthFt
}

// TimeBin implements Equation (5):
//
//	tb_i = (t_i // 10) * 10
//
// Discretizes continuous time (seconds) into discrete temporal windows.
func TimeBin(tS float64, binWidthS int) int {
	if binWidthS <= 0 {
		binWidthS = TimeBinS
	}
	return int(math.Floor(tS/float64(binWidthS))) * binWidthS
}

// ToVoxelKey converts 4D coordinates (latitude, longitude, altitude_ft, time_seconds)
// into a composite 4D VoxelKey (Eq. 2 & 6).
func ToVoxelKey(lat, lon, altFt, tS float64, resolution int) (VoxelKey, error) {
	cell, err := spatial.PointToH3Cell(lat, lon, resolution)
	if err != nil {
		return VoxelKey{}, fmt.Errorf("mapping point to H3 cell: %w", err)
	}

	return VoxelKey{
		H3Cell:   cell,
		AltBinFt: AltitudeBin(altFt, AltitudeBinFt),
		TimeBinS: TimeBin(tS, TimeBinS),
	}, nil
}
