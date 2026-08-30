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

// Package protectedvolumes models static and dynamic restricted airspace volumes (vertiports, geofenced corridors).
// Populates the 4D occupancy map with RESERVED:: markers so that entering aircraft register instant conflicts
// without adding special-case geometry paths to the core detection engine.
package protectedvolumes

import (
	"context"

	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/occupancy"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/spatial"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/temporal"
)

// ReservedEntityPrefix is the sentinel prefix identifying protected airspace constraints.
const ReservedEntityPrefix = "RESERVED::"

// ProtectedVolume represents a 3D cylindrical/hexagonal geofence volume across altitude bands.
type ProtectedVolume struct {
	VolumeID      string
	CenterLat     float64
	CenterLon     float64
	RadiusH3Rings int
	FloorFt       float64
	CeilingFt     float64
}

// VoxelKeysForTimeBin generates all composite 4D voxel keys occupied by this protected volume during a given time bin.
func (v ProtectedVolume) VoxelKeysForTimeBin(timeBinS, resolution int) ([]temporal.VoxelKey, error) {
	centerCell, err := spatial.PointToH3Cell(v.CenterLat, v.CenterLon, resolution)
	if err != nil {
		return nil, err
	}

	cells, err := spatial.NeighborCells(centerCell, v.RadiusH3Rings)
	if err != nil {
		return nil, err
	}

	var keys []temporal.VoxelKey
	for alt := v.FloorFt; alt <= v.CeilingFt; alt += float64(temporal.AltitudeBinFt) {
		ab := temporal.AltitudeBin(alt, temporal.AltitudeBinFt)
		for _, c := range cells {
			keys = append(keys, temporal.VoxelKey{
				H3Cell:   c,
				AltBinFt: ab,
				TimeBinS: timeBinS,
			})
		}
	}
	return keys, nil
}

// Register pre-populates the shared occupancy map with constraint entities.
func Register(ctx context.Context, occ occupancy.Map, v ProtectedVolume, timeBinS, resolution int) error {
	keys, err := v.VoxelKeysForTimeBin(timeBinS, resolution)
	if err != nil {
		return err
	}

	reservedID := ReservedEntityPrefix + v.VolumeID
	for _, key := range keys {
		if err := occ.Add(ctx, key, reservedID); err != nil {
			return err
		}
	}
	return nil
}
