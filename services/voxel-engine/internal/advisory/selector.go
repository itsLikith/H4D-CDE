package advisory

import (
	"context"
	"strconv"

	"github.com/uber/h3-go/v4"
)

var DepartureDelayOptionsS = []float64{30, 60, 90, 120}

const AltitudeOffsetFt = 500

// Expected risk reductions, as reported in the paper's discussion of Table III.
const (
	RiskReductionReroute        = 0.90
	RiskReductionAltitudeChange = 0.85
	RiskReductionDelay          = 0.70
)

type Strategy string

const (
	StrategyDelay          Strategy = "delay"
	StrategyReroute        Strategy = "reroute"
	StrategyAltitudeChange Strategy = "altitude_change"
)

type Result struct {
	Strategy                 Strategy
	Parameters               map[string]string
	ExpectedRiskReductionPct float64
}

type Conflict struct {
	ConflictID           string
	OriginCell, DestCell h3.Cell
}

// ConflictResolutionChecker re-runs the conflict check with a hypothetical change applied, without committing that change.
type ConflictResolutionChecker interface {
	ResolvesWithDelay(ctx context.Context, c Conflict, delayS float64) (bool, error)
	ResolvesWithPath(ctx context.Context, c Conflict, path []h3.Cell) (bool, error)
}

// SelectAdvisory implements the paper's greedy cascade :
// try delay, then a rerouted path, then an altitude change — stopping at the first option that actually clears the conflict.
func SelectAdvisory(ctx context.Context, c Conflict, checker ConflictResolutionChecker, w Weights, k int) (Result, error) {
	for _, delayS := range DepartureDelayOptionsS {
		ok, err := checker.ResolvesWithDelay(ctx, c, delayS)
		if err != nil {
			return Result{}, err
		}
		if ok {
			return Result{
				Strategy:                 StrategyDelay,
				Parameters:               map[string]string{"delay_s": strconv.FormatFloat(delayS, 'f', 0, 64)},
				ExpectedRiskReductionPct: RiskReductionDelay * 100,
			}, nil
		}
	}

	for _, path := range kShortestPathReroute(c.OriginCell, c.DestCell, k, 20) {
		resolved, err := checker.ResolvesWithPath(ctx, c, path)
		if err == nil && resolved {
			return Result{
				Strategy:                 StrategyReroute,
				Parameters:               map[string]string{"path_hops": strconv.Itoa(len(path) - 1)},
				ExpectedRiskReductionPct: RiskReductionReroute * 100,
			}, nil
		}
	}

	// Altitude change is the cascade's guaranteed fallback -- it never
	// itself fails, unlike the delay/reroute branches above which can
	// legitimately fail to resolve anything.
	return Result{
		Strategy:                 StrategyAltitudeChange,
		Parameters:               map[string]string{"delta_ft": strconv.Itoa(AltitudeOffsetFt)},
		ExpectedRiskReductionPct: RiskReductionAltitudeChange * 100,
	}, nil
}
