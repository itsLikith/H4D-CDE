// services/voxel-engine/internal/conflict/conflict.go
package conflict

import (
	"context"
	"math"

	"hive/voxel-engine/internal/occupancy"
	"hive/voxel-engine/internal/spatial"
	"hive/voxel-engine/internal/temporal"
)

const earthRadiusKm = 6371.0088

// HaversineNM returns great-circle distance in nautical miles.
func HaversineNM(lat1, lon1, lat2, lon2 float64) float64 {
	p1, p2 := lat1*math.Pi/180, lat2*math.Pi/180
	dphi, dlmb := (lat2-lat1)*math.Pi/180, (lon2-lon1)*math.Pi/180
	a := math.Sin(dphi/2)*math.Sin(dphi/2) + math.Cos(p1)*math.Cos(p2)*math.Sin(dlmb/2)*math.Sin(dlmb/2)
	return 2 * earthRadiusKm * math.Asin(math.Sqrt(a)) * 0.539957 // km -> NM
}

type Position struct{ Lat, Lon, AltFt float64 }
type Pair struct {
	EntityA, EntityB string
	Key              temporal.VoxelKey
}

// SameVoxelConflicts: stage (A) -- O(n) over occupied keys.
func SameVoxelConflicts(ctx context.Context, occ occupancy.Map, keys []temporal.VoxelKey) ([]Pair, error) {
	var out []Pair
	for _, key := range keys {
		occupants, err := occ.Occupants(ctx, key)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(occupants); i++ {
			for j := i + 1; j < len(occupants); j++ {
				out = append(out, Pair{occupants[i], occupants[j], key})
			}
		}
	}
	return out, nil
}

// NeighborVoxelConflicts: stage (B) -- O(n) over occupied keys, with a constant factor of 6*3 = 18
// Exactly 6 neighbours x 3 altitude bins
// per occupied voxel: a constant amount of work regardless of total traffic,
// which is why this stays O(n) overall.
func NeighborVoxelConflicts(
	ctx context.Context, occ occupancy.Map, keys []temporal.VoxelKey,
	positions map[string]Position, hSepNM, vSepFt float64,
) ([]Pair, error) {
	seen := map[string]bool{}
	var out []Pair

	for _, key := range keys {
		ownOccupants, err := occ.Occupants(ctx, key)
		if err != nil {
			return nil, err
		}
		neighbors, err := spatial.NeighborCells(key.H3Cell, 1)
		if err != nil {
			return nil, err
		}
		for _, nbrCell := range neighbors {
			if nbrCell == key.H3Cell {
				continue // that's the same-voxel case, handled above
			}
			for _, altDelta := range []int{-temporal.AltitudeBinFt, 0, temporal.AltitudeBinFt} {
				nbrKey := temporal.VoxelKey{H3Cell: nbrCell, AltBinFt: key.AltBinFt + altDelta, TimeBinS: key.TimeBinS}
				nbrOccupants, err := occ.Occupants(ctx, nbrKey)
				if err != nil {
					return nil, err
				}
				for _, e1 := range ownOccupants {
					for _, e2 := range nbrOccupants {
						if e1 == e2 {
							continue
						}
						pid := pairID(e1, e2)
						if seen[pid] {
							continue
						}
						p1, ok1 := positions[e1]
						p2, ok2 := positions[e2]
						if !ok1 || !ok2 {
							continue
						}
						if HaversineNM(p1.Lat, p1.Lon, p2.Lat, p2.Lon) < hSepNM && math.Abs(p1.AltFt-p2.AltFt) < vSepFt {
							seen[pid] = true
							out = append(out, Pair{e1, e2, nbrKey})
						}
					}
				}
			}
		}
	}
	return out, nil
}

func pairID(a, b string) string {
	if a < b {
		return a + "|" + b
	}
	return b + "|" + a
}
