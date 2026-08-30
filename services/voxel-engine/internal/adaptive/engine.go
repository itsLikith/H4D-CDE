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

// Package adaptive implements the Adaptive Discretization Engine (Module 6).
// Fulfills Research Objective 2 from the paper by dynamically refining spatial H3 resolution (res 8 -> res 9)
// and temporal binning (10s -> 5s) in localized high-density airspace sectors to prevent false negatives.
package adaptive

import "github.com/uber/h3-go/v4"

type Config struct {
	BaseResolution   int
	FineResolution   int
	BaseTimeBinS     int
	FineTimeBinS     int
	DensityThreshold int // Occupancy forecast threshold that triggers high-resolution zoom
}

// DefaultConfig provides standard operational thresholds for adaptive discretization.
func DefaultConfig() Config {
	return Config{
		BaseResolution:   8,
		FineResolution:   9,
		BaseTimeBinS:     10,
		FineTimeBinS:     5,
		DensityThreshold: 6,
	}
}

type Engine struct {
	cfg Config
}

func New(cfg Config) *Engine {
	return &Engine{cfg: cfg}
}

// ResolutionFor selects the optimal H3 resolution given the forecasted local sector occupancy.
func (e *Engine) ResolutionFor(forecastOccupancy int) int {
	if forecastOccupancy >= e.cfg.DensityThreshold {
		return e.cfg.FineResolution
	}
	return e.cfg.BaseResolution
}

// TimeBinFor selects the time bin width in seconds.
func (e *Engine) TimeBinFor(forecastOccupancy int) int {
	if forecastOccupancy >= e.cfg.DensityThreshold {
		return e.cfg.FineTimeBinS
	}
	return e.cfg.BaseTimeBinS
}

// RefineCell subdivides a parent H3 cell into its 7 finer child hexagons at FineResolution.
func (e *Engine) RefineCell(parent h3.Cell) ([]h3.Cell, error) {
	return parent.Children(e.cfg.FineResolution)
}
