// services/standards-svc/internal/icaofpl/parser.go
package icaofpl

import (
	"fmt"

	flightplanv1 "hive/gen/flightplan"
)

// ValidationError collects every problem found, rather than failing on the
// first one -- an operator resubmitting a flight plan wants the whole list
// of what to fix, not one error at a time.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string { return fmt.Sprintf("%d validation issue(s)", len(e.Issues)) }

// Validate enforces the plausibility checks flagged as a security control
// rejecting malformed or implausible submissions here, at
// the system's edge, before they ever reach voxel-engine.
func Validate(fpl *flightplanv1.FlightPlan) error {
	var issues []string
	if len(fpl.Waypoints) < 2 {
		issues = append(issues, "flight plan needs at least an origin and destination waypoint")
	}
	if fpl.CruiseAltitudeFt <= 0 || fpl.CruiseAltitudeFt > 60000 {
		issues = append(issues, "cruise altitude out of plausible range")
	}
	if fpl.CruiseSpeedKt <= 0 || fpl.CruiseSpeedKt > 400 {
		issues = append(issues, "cruise speed out of plausible range for a UAM/UAS airframe")
	}
	for _, wp := range fpl.Waypoints {
		if wp.Lat < -90 || wp.Lat > 90 || wp.Lon < -180 || wp.Lon > 180 {
			issues = append(issues, "waypoint coordinates out of range")
		}
	}
	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}