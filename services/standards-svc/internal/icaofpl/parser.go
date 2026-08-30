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

// Package icaofpl validates and parses ICAO Doc 4444 Flight Plans (FPL) and ASTM F3548-21 Operational Intent.
package icaofpl

import (
	"fmt"
	"strings"

	flightplanv1 "github.com/itsLikith/h4d-cde/gen/flightplan"
)

// ValidationError collects all validation issues so an operator can fix all problems in one submission.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%d validation issue(s): %s", len(e.Issues), strings.Join(e.Issues, "; "))
}

// Validate executes strict plausibility and airspace boundary checks on submitted flight plans.
func Validate(fpl *flightplanv1.FlightPlan) error {
	if fpl == nil {
		return &ValidationError{Issues: []string{"flight plan cannot be nil"}}
	}

	var issues []string

	if strings.TrimSpace(fpl.EntityId) == "" {
		issues = append(issues, "entity_id (aircraft identification) is required")
	}

	if len(fpl.Waypoints) < 2 {
		issues = append(issues, "flight plan requires at least origin and destination waypoints")
	}

	if fpl.CruiseAltitudeFt <= 0 || fpl.CruiseAltitudeFt > 60000 {
		issues = append(issues, "cruise altitude out of valid range (0, 60000] ft")
	}

	if fpl.CruiseSpeedKt <= 0 || fpl.CruiseSpeedKt > 450 {
		issues = append(issues, "cruise speed out of valid range (0, 450] kt for UAM/UAS airframes")
	}

	for i, wp := range fpl.Waypoints {
		if wp.Lat < -90.0 || wp.Lat > 90.0 || wp.Lon < -180.0 || wp.Lon > 180.0 {
			issues = append(issues, fmt.Sprintf("waypoint #%d coordinates (%.4f, %.4f) out of geographic bounds", i+1, wp.Lat, wp.Lon))
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}

	return nil
}