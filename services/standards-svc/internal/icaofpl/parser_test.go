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

package icaofpl_test

import (
	"testing"

	commonv1 "github.com/itsLikith/h4d-cde/gen/common"
	flightplanv1 "github.com/itsLikith/h4d-cde/gen/flightplan"
	"github.com/itsLikith/h4d-cde/services/standards-svc/internal/icaofpl"
	"github.com/stretchr/testify/assert"
)

func TestValidFlightPlan(t *testing.T) {
	fpl := &flightplanv1.FlightPlan{
		EntityId:        "UAV-VALID-01",
		OriginIcao:      "OMDB",
		DestinationIcao: "OMAA",
		EobtUnixMs:      1700000000000,
		CruiseAltitudeFt: 1500.0,
		CruiseSpeedKt:   90.0,
		Waypoints: []*commonv1.GeoPoint{
			{Lat: 25.2532, Lon: 55.3657},
			{Lat: 24.4330, Lon: 54.6511},
		},
	}

	assert.NoError(t, icaofpl.Validate(fpl))
}

func TestInvalidFlightPlanRejections(t *testing.T) {
	// Missing waypoints & invalid speed/altitude
	fpl := &flightplanv1.FlightPlan{
		EntityId:        "",
		CruiseAltitudeFt: -100,
		CruiseSpeedKt:   1200,
		Waypoints:       []*commonv1.GeoPoint{},
	}

	err := icaofpl.Validate(fpl)
	assert.Error(t, err)
	valErr, ok := err.(*icaofpl.ValidationError)
	assert.True(t, ok)
	assert.GreaterOrEqual(t, len(valErr.Issues), 3)
}
