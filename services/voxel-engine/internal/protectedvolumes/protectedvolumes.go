package protectedvolumes

import (
	"context"

	"hive/voxel-engine/internal/occupancy"
	"hive/voxel-engine/internal/spatial"
	"hive/voxel-engine/internal/temporal"
)

const ReservedEntityPrefix = "RESERVED::"

type ProtectedVolume struct {
	VolumeID      string
	CenterLat     float64
	CenterLon     float64
	RadiusH3Rings int
	FloorFt       float64
	CeilingFt     float64
}

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
	for alt := v.FloorFt; alt <= v.CeilingFt; alt += temporal.AltitudeBinFt {
		ab := temporal.AltitudeBin(alt, temporal.AltitudeBinFt)
		for _, c := range cells {
			keys = append(keys, temporal.VoxelKey{H3Cell: c, AltBinFt: ab, TimeBinS: timeBinS})
		}
	}
	return keys, nil
}

// Register pre-populates the occupancy map so any real aircraft entering
// these voxels immediately registers as a conflict against the reserved
// marker  without special-casing
// protected volumes anywhere in the conflict-check logic itself.
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
