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

package advisory

import (
	"context"
	"fmt"
	"strconv"

	"github.com/uber/h3-go/v4"
)

// Standard departure delay increments tested in the cascade.
var DepartureDelayOptionsS = []float64{30, 60, 90, 120}

// AltitudeOffsetFt is the vertical flight level adjustment (±500 ft).
const AltitudeOffsetFt = 500

// Theoretical risk reductions reported in Table III of Sahadevan et al. (ICSPIS 2025):
const (
	RiskReductionReroute        = 0.90 // 90% expected risk reduction
	RiskReductionAltitudeChange = 0.85 // 85% expected risk reduction
	RiskReductionDelay          = 0.70 // 70% expected risk reduction
)

type Strategy string

const (
	StrategyDelay          Strategy = "delay"
	StrategyReroute        Strategy = "reroute"
	StrategyAltitudeChange Strategy = "altitude_change"
)

// Result represents the optimal synthesized advisory recommended for an operator.
type Result struct {
	ConflictID               string
	Strategy                 Strategy
	Parameters               map[string]string
	ExpectedRiskReductionPct float64
}

type Conflict struct {
	ConflictID           string
	EntityA              string
	EntityB              string
	OriginCell, DestCell h3.Cell
}

// ConflictResolutionChecker checks whether a candidate modification resolves the spatial/temporal overlap.
type ConflictResolutionChecker interface {
	ResolvesWithDelay(ctx context.Context, c Conflict, delayS float64) (bool, error)
	ResolvesWithPath(ctx context.Context, c Conflict, path []h3.Cell) (bool, error)
}

// SelectAdvisory executes the hierarchical greedy resolution cascade:
//  1. Lateral / Departure Delay: least disruptive to inflight profiles.
//  2. Spatial Reroute: alternative H3 corridor paths avoiding congested voxels.
//  3. Vertical Flight Level Change: guaranteed fallback separation (±500 ft offset).
func SelectAdvisory(
	ctx context.Context,
	c Conflict,
	checker ConflictResolutionChecker,
	w Weights,
	kAlternativePaths int,
) (Result, error) {
	// 1. Evaluate temporal ground delay
	for _, delayS := range DepartureDelayOptionsS {
		if checker != nil {
			resolved, err := checker.ResolvesWithDelay(ctx, c, delayS)
			if err == nil && resolved {
				return Result{
					ConflictID:               c.ConflictID,
					Strategy:                 StrategyDelay,
					Parameters:               map[string]string{"delay_s": strconv.FormatFloat(delayS, 'f', 0, 64)},
					ExpectedRiskReductionPct: RiskReductionDelay * 100.0,
				}, nil
			}
		} else {
			// Fast path without mock checker
			return Result{
				ConflictID:               c.ConflictID,
				Strategy:                 StrategyDelay,
				Parameters:               map[string]string{"delay_s": fmt.Sprintf("%.0f", delayS)},
				ExpectedRiskReductionPct: RiskReductionDelay * 100.0,
			}, nil
		}
	}

	// 2. Evaluate lateral rerouting via H3 neighbor graph
	if c.OriginCell != 0 && c.DestCell != 0 {
		altPaths := kShortestPathReroute(c.OriginCell, c.DestCell, kAlternativePaths, 20)
		for _, path := range altPaths {
			if checker != nil {
				resolved, err := checker.ResolvesWithPath(ctx, c, path)
				if err == nil && resolved {
					return Result{
						ConflictID:               c.ConflictID,
						Strategy:                 StrategyReroute,
						Parameters:               map[string]string{"path_hops": strconv.Itoa(len(path) - 1)},
						ExpectedRiskReductionPct: RiskReductionReroute * 100.0,
					}, nil
				}
			}
		}
	}

	// 3. Fallback: Vertical separation altitude change
	return Result{
		ConflictID:               c.ConflictID,
		Strategy:                 StrategyAltitudeChange,
		Parameters:               map[string]string{"delta_ft": strconv.Itoa(AltitudeOffsetFt)},
		ExpectedRiskReductionPct: RiskReductionAltitudeChange * 100.0,
	}, nil
}
