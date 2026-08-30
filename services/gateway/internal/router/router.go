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

// Package router provides the REST API endpoints, gRPC client dispatching, and WebSocket upgrade routes.
package router

import (
	"strconv"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	auditv1 "github.com/itsLikith/h4d-cde/gen/audit"
	commonv1 "github.com/itsLikith/h4d-cde/gen/common"
	flightplanv1 "github.com/itsLikith/h4d-cde/gen/flightplan"
	standardsv1 "github.com/itsLikith/h4d-cde/gen/standards"
	voxelizerv1 "github.com/itsLikith/h4d-cde/gen/voxelizer"

	"github.com/itsLikith/h4d-cde/services/gateway/internal/auth"
	"github.com/itsLikith/h4d-cde/services/gateway/internal/ws"
)

type Clients struct {
	StandardsClient   standardsv1.StandardsServiceClient
	VoxelEngineClient voxelizerv1.VoxelEngineServiceClient
	AuditClient       auditv1.AuditServiceClient
}

type GeoPointJSON struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type FlightPlanRequest struct {
	EntityID         string         `json:"entity_id"`
	OriginICAO       string         `json:"origin_icao"`
	DestinationICAO  string         `json:"destination_icao"`
	EOBTUnixMs       int64          `json:"eobt_unix_ms"`
	Waypoints        []GeoPointJSON `json:"waypoints"`
	CruiseAltitudeFt float64        `json:"cruise_altitude_ft"`
	CruiseSpeedKt    float64        `json:"cruise_speed_kt"`
}

func Register(app *fiber.App, hub *ws.Hub, clients Clients) {
	// Public Health & WebSocket
	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "gateway",
			"engine":  "H4D-CDE",
		})
	})

	app.Get("/ws/live-updates", websocket.New(func(c *websocket.Conn) {
		hub.Register(c)
		defer hub.Unregister(c)
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))

	// Protected v1 API Group
	v1 := app.Group("/v1", auth.RequireToken())

	v1.Post("/flight-plans", func(c fiber.Ctx) error {
		var req FlightPlanRequest
		if err := c.Bind().Body(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if clients.StandardsClient == nil {
			return fiber.NewError(fiber.StatusBadGateway, "standards-svc unavailable")
		}

		var protoWaypoints []*commonv1.GeoPoint
		for _, wp := range req.Waypoints {
			protoWaypoints = append(protoWaypoints, &commonv1.GeoPoint{
				Lat: wp.Lat,
				Lon: wp.Lon,
			})
		}

		protoFPL := &flightplanv1.FlightPlan{
			EntityId:         req.EntityID,
			OriginIcao:       req.OriginICAO,
			DestinationIcao:  req.DestinationICAO,
			EobtUnixMs:       req.EOBTUnixMs,
			Waypoints:        protoWaypoints,
			CruiseAltitudeFt: req.CruiseAltitudeFt,
			CruiseSpeedKt:    req.CruiseSpeedKt,
		}

		result, err := clients.StandardsClient.SubmitFlightPlan(c.Context(), &standardsv1.SubmitFlightPlanRequest{
			FlightPlan: protoFPL,
		})
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.Status(fiber.StatusAccepted).JSON(result)
	})

	v1.Get("/conflicts", func(c fiber.Ctx) error {
		if clients.VoxelEngineClient == nil {
			return fiber.NewError(fiber.StatusBadGateway, "voxel-engine unavailable")
		}

		h3Cell := c.Query("h3_cell")
		minRiskStr := c.Query("min_risk")
		minRisk := 0.0
		if minRiskStr != "" {
			minRisk, _ = strconv.ParseFloat(minRiskStr, 64)
		}

		resp, err := clients.VoxelEngineClient.GetConflicts(c.Context(), &voxelizerv1.GetConflictsRequest{
			H3Cell:  h3Cell,
			MinRisk: minRisk,
		})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(resp.Conflicts)
	})

	v1.Get("/audit/voxel/:h3Cell/:altBin/:timeBin", func(c fiber.Ctx) error {
		if clients.AuditClient == nil {
			return fiber.NewError(fiber.StatusBadGateway, "audit-svc unavailable")
		}

		h3Cell := c.Params("h3Cell")
		altBin, _ := strconv.Atoi(c.Params("altBin"))
		timeBin, _ := strconv.Atoi(c.Params("timeBin"))

		resp, err := clients.AuditClient.GetVoxelAuditTrail(c.Context(), &auditv1.GetVoxelAuditTrailRequest{
			VoxelKey: &commonv1.VoxelKey{
				H3Cell:   h3Cell,
				AltBinFt: int32(altBin),
				TimeBinS: int32(timeBin),
			},
		})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(resp)
	})
}