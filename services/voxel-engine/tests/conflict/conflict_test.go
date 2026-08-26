package conflict_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"hive/voxel-engine/internal/conflict"
	"hive/voxel-engine/internal/occupancy"
	"hive/voxel-engine/internal/temporal"
)

func TestSameVoxelConflictDetected(t *testing.T) {
	ctx := context.Background()
	occ := occupancy.NewInMemoryMap()
	key, err := temporal.ToVoxelKey(25.20, 55.27, 480, 10004, 8)
	assert.NoError(t, err)
	assert.NoError(t, occ.Add(ctx, key, "UAV-A"))
	assert.NoError(t, occ.Add(ctx, key, "UAV-B"))

	pairs, err := conflict.SameVoxelConflicts(ctx, occ, []temporal.VoxelKey{key})
	assert.NoError(t, err)
	assert.Len(t, pairs, 1)
}

func TestNoConflictForSingleOccupant(t *testing.T) {
	ctx := context.Background()
	occ := occupancy.NewInMemoryMap()
	key, _ := temporal.ToVoxelKey(25.20, 55.27, 480, 10004, 8)
	assert.NoError(t, occ.Add(ctx, key, "UAV-A"))

	pairs, err := conflict.SameVoxelConflicts(ctx, occ, []temporal.VoxelKey{key})
	assert.NoError(t, err)
	assert.Empty(t, pairs)
}

func TestAltitudeBinFloorsNotRounds(t *testing.T) {
	assert.Equal(t, 400, temporal.AltitudeBin(499, temporal.AltitudeBinFt))
	assert.Equal(t, 500, temporal.AltitudeBin(500, temporal.AltitudeBinFt))
}