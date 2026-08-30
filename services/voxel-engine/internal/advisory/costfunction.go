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

// Package advisory implements the Strategic Conflict Resolution and Advisory Selector (Module 5).
package advisory

// Weights represents the multi-objective cost coefficients (w1, w2, w3, w4) from Equation (11):
//
//	C = w1 * delay + w2 * path_dev + w3 * alt_change + w4 * conflict_prob
type Weights struct {
	Delay               float64
	PathDeviation       float64
	AltitudeChange      float64
	ConflictProbability float64
}

// DefaultWeights initializes standard operational cost coefficients.
func DefaultWeights() Weights {
	return Weights{
		Delay:               0.30,
		PathDeviation:       0.30,
		AltitudeChange:      0.20,
		ConflictProbability: 0.20,
	}
}

// Cost computes the aggregate operational penalty for an advisory candidate (Eq. 11).
func Cost(delayS, pathDeviationKm, altitudeChangeFt, conflictProbability float64, w Weights) float64 {
	return w.Delay*delayS +
		w.PathDeviation*pathDeviationKm +
		w.AltitudeChange*altitudeChangeFt +
		w.ConflictProbability*conflictProbability
}
