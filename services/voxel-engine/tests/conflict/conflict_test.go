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

package conflict_test

import (
	"context"
	"testing"

	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/conflict"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/occupancy"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/temporal"
	"github.com/stretchr/testify/assert"
)

func TestSameVoxelConflictDetected(t *testing.T) {
	ctx := context.Background()
	occ := occupancy.NewInMemoryMap()

	// Two UAVs landing in the exact same 4D voxel bucket
	key, err := temporal.ToVoxelKey(25.20, 55.27, 480, 10004, 8)
	assert.NoError(t, err)

	assert.NoError(t, occ.Add(ctx, key, "UAV-A"))
	assert.NoError(t, occ.Add(ctx, key, "UAV-B"))

	pairs, err := conflict.SameVoxelConflicts(ctx, occ, []temporal.VoxelKey{key})
	assert.NoError(t, err)
	assert.Len(t, pairs, 1)
	assert.Equal(t, conflict.ConflictTypeSameVoxel, pairs[0].ConflictType)
	assert.ElementsMatch(t, []string{"UAV-A", "UAV-B"}, []string{pairs[0].EntityA, pairs[0].EntityB})
}

func TestNoConflictForSingleOccupant(t *testing.T) {
	ctx := context.Background()
	occ := occupancy.NewInMemoryMap()

	key, err := temporal.ToVoxelKey(25.20, 55.27, 480, 10004, 8)
	assert.NoError(t, err)
	assert.NoError(t, occ.Add(ctx, key, "UAV-A"))

	pairs, err := conflict.SameVoxelConflicts(ctx, occ, []temporal.VoxelKey{key})
	assert.NoError(t, err)
	assert.Empty(t, pairs)
}

func TestAltitudeBinFloorsAccurately(t *testing.T) {
	assert.Equal(t, 400, temporal.AltitudeBin(499, temporal.AltitudeBinFt))
	assert.Equal(t, 500, temporal.AltitudeBin(500, temporal.AltitudeBinFt))
	assert.Equal(t, 0, temporal.AltitudeBin(0, temporal.AltitudeBinFt))
}

func TestTimeBinFloorsAccurately(t *testing.T) {
	assert.Equal(t, 10, temporal.TimeBin(14, temporal.TimeBinS))
	assert.Equal(t, 20, temporal.TimeBin(20, temporal.TimeBinS))
}
