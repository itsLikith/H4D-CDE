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

	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/adaptive"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/advisory"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/conflict"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/occupancy"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/protectedvolumes"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/spatial"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/temporal"
	"github.com/stretchr/testify/assert"
)

func TestNeighborVoxelConflictDetected(t *testing.T) {
	ctx := context.Background()
	occ := occupancy.NewInMemoryMap()

	lat1, lon1, alt1 := 25.2532, 55.3657, 1000.0
	key1, err := temporal.ToVoxelKey(lat1, lon1, alt1, 100, 8)
	assert.NoError(t, err)

	// Pick an immediate adjacent neighbor cell
	nbrs, err := spatial.NeighborCells(key1.H3Cell, 1)
	assert.NoError(t, err)
	assert.Len(t, nbrs, 7)

	nbrCell := nbrs[1]
	key2 := temporal.VoxelKey{
		H3Cell:   nbrCell,
		AltBinFt: 1000,
		TimeBinS: 100,
	}

	assert.NoError(t, occ.Add(ctx, key1, "UAV-1"))
	assert.NoError(t, occ.Add(ctx, key2, "UAV-2"))

	latLng2, err := nbrCell.LatLng()
	assert.NoError(t, err)

	positions := map[string]conflict.Position{
		"UAV-1": {Lat: lat1, Lon: lon1, AltFt: alt1},
		"UAV-2": {Lat: latLng2.Lat, Lon: latLng2.Lng, AltFt: 1050.0},
	}

	keys := []temporal.VoxelKey{key1, key2}
	nbrConflicts, err := conflict.NeighborVoxelConflicts(ctx, occ, keys, positions, 5.0, 1000.0)
	assert.NoError(t, err)
	assert.NotEmpty(t, nbrConflicts)
	assert.Equal(t, conflict.ConflictTypeNeighborVoxel, nbrConflicts[0].ConflictType)
}

func TestAdvisoryCostFunction(t *testing.T) {
	weights := advisory.DefaultWeights()
	// Cost = 0.3*60 + 0.3*2.0 + 0.2*500 + 0.2*0.8 = 18 + 0.6 + 100 + 0.16 = 118.76
	cost := advisory.Cost(60.0, 2.0, 500.0, 0.8, weights)
	assert.InDelta(t, 118.76, cost, 0.01)
}

func TestAdvisoryCascadeSelectsDelayOrAltitude(t *testing.T) {
	ctx := context.Background()
	c := advisory.Conflict{
		ConflictID: "CONF-TEST-001",
		EntityA:    "UAV-1",
		EntityB:    "UAV-2",
	}

	res, err := advisory.SelectAdvisory(ctx, c, nil, advisory.DefaultWeights(), 3)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Strategy)
	assert.Greater(t, res.ExpectedRiskReductionPct, 0.0)
}

func TestProtectedVolumeReservation(t *testing.T) {
	ctx := context.Background()
	occ := occupancy.NewInMemoryMap()

	pv := protectedvolumes.ProtectedVolume{
		VolumeID:      "VERTIPORT-DXB",
		CenterLat:     25.2532,
		CenterLon:     55.3657,
		RadiusH3Rings: 1,
		FloorFt:       0,
		CeilingFt:     500,
	}

	timeBinS := 10000
	assert.NoError(t, protectedvolumes.Register(ctx, occ, pv, timeBinS, 8))

	// Aircraft entering the vertiport volume
	acKey, err := temporal.ToVoxelKey(25.2532, 55.3657, 300, float64(timeBinS), 8)
	assert.NoError(t, err)
	assert.NoError(t, occ.Add(ctx, acKey, "UAV-TRESPASSER"))

	pairs, err := conflict.SameVoxelConflicts(ctx, occ, []temporal.VoxelKey{acKey})
	assert.NoError(t, err)
	assert.NotEmpty(t, pairs)
	assert.Contains(t, []string{pairs[0].EntityA, pairs[0].EntityB}, "RESERVED::VERTIPORT-DXB")
	assert.Contains(t, []string{pairs[0].EntityA, pairs[0].EntityB}, "UAV-TRESPASSER")
}

func TestAdaptiveDiscretizationRefinement(t *testing.T) {
	engine := adaptive.New(adaptive.DefaultConfig())

	// Low density (< 6) -> Base res 8
	assert.Equal(t, 8, engine.ResolutionFor(3))
	assert.Equal(t, 10, engine.TimeBinFor(3))

	// High density (>= 6) -> Fine res 9, 5s time bin
	assert.Equal(t, 9, engine.ResolutionFor(8))
	assert.Equal(t, 5, engine.TimeBinFor(8))

	cell, err := spatial.PointToH3Cell(25.25, 55.30, 8)
	assert.NoError(t, err)
	children, err := engine.RefineCell(cell)
	assert.NoError(t, err)
	assert.Len(t, children, 7)
}
