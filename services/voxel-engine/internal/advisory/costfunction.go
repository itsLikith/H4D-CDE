package advisory

// Weights are Eq. (11)'s w1..w4 -- config, not constants.
type Weights struct {
	Delay               float64
	PathDeviation       float64
	AltitudeChange      float64
	ConflictProbability float64
}

func DefaultWeights() Weights {
	return Weights{Delay: 0.30, PathDeviation: 0.30, AltitudeChange: 0.20, ConflictProbability: 0.20}
}

// Cost implements Eq. (11): C = w1*delay + w2*path_dev + w3*alt_change + w4*conflict_prob.
func Cost(delayS, pathDeviationKm, altitudeChangeFt, conflictProbability float64, w Weights) float64 {
	return w.Delay*delayS + w.PathDeviation*pathDeviationKm +
		w.AltitudeChange*altitudeChangeFt + w.ConflictProbability*conflictProbability
}
