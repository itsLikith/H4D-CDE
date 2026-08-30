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

// Package baseline implements the unoptimized, quadratic O(n²) pairwise conflict detection baseline
// benchmarked in Section II.C and Table I of Sahadevan et al. (ICSPIS 2025).
package baseline

import (
	"math"
)

type Point struct {
	EntityID string
	TS       float64
	Lat      float64
	Lon      float64
	AltFt    float64
}

const earthRadiusKm = 6371.0088

func haversineNM(lat1, lon1, lat2, lon2 float64) float64 {
	p1, p2 := lat1*math.Pi/180.0, lat2*math.Pi/180.0
	dphi := (lat2 - lat1) * math.Pi / 180.0
	dlmb := (lon2 - lon1) * math.Pi / 180.0

	a := math.Sin(dphi/2.0)*math.Sin(dphi/2.0) +
		math.Cos(p1)*math.Cos(p2)*math.Sin(dlmb/2.0)*math.Sin(dlmb/2.0)
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))
	return earthRadiusKm * c * 0.539957 // Convert km to NM
}

// PairwiseConflictDetection executes exhaustive O(n²) comparison across all pairs of trajectory points.
// Checks temporal overlap (|t1 - t2| <= tSepS), horizontal separation (< hSepNM), and vertical separation (< vSepFt).
func PairwiseConflictDetection(points []Point, hSepNM, vSepFt, tSepS float64) [][2]Point {
	var conflicts [][2]Point
	n := len(points)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			a, b := points[i], points[j]

			// Skip same entity comparisons
			if a.EntityID == b.EntityID {
				continue
			}

			// Temporal window filter
			if math.Abs(a.TS-b.TS) > tSepS {
				continue
			}

			// Great-circle horizontal and vertical distance comparison
			if haversineNM(a.Lat, a.Lon, b.Lat, b.Lon) < hSepNM && math.Abs(a.AltFt-b.AltFt) < vSepFt {
				conflicts = append(conflicts, [2]Point{a, b})
			}
		}
	}

	return conflicts
}
