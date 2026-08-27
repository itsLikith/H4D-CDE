package temporal

import (
	"fmt"
	"hive/voxel-engine/internal/spatial"

	"github.com/uber/h3-go/v4"
)

const (
	AltitudeBinFt = 100
	TimeBinS      = 10
)

// VoxelKey is Eq.: voxel_key = (h_voxel, a_bin, t_bin).
type VoxelKey struct {
	H3Cell   h3.Cell
	AltBinFt int
	TimeBinS int
}

// String renders the key flat, e.g. "882830...:400:36000" -- used as the
// Redis key as the wire representation in proto/common.proto.
func (k VoxelKey) String() string {
	return fmt.Sprintf("%s:%d:%d", k.H3Cell.String(), k.AltBinFt, k.TimeBinS)
}

// AltitudeBin implements Eq.: ab = (alt // 100) * 100.
func AltitudeBin(altFt float64, binWidthFt int) int {
	return int(altFt) / binWidthFt * binWidthFt
}

// TimeBin implements Eq.: tb = (t // 10) * 10.
func TimeBin(tS float64, binWidthS int) int {
	return int(tS) / binWidthS * binWidthS
}

// ToVoxelKey implements Eq.: voxel_key = (h_voxel, a_bin, t_bin).
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
