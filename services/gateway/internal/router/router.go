// services/gateway/internal/router/router.go
package router

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v3"
	"hive/gateway/internal/auth"
	"hive/gateway/internal/ws"
)

func Register(app *fiber.App, hub *ws.Hub) {
	v1 := app.Group("/v1", auth.RequireToken())

	v1.Post("/flight-plans", submitFlightPlan)
	v1.Get("/conflicts", listConflicts)
	v1.Get("/conflicts/:id", getConflict)
	v1.Post("/conflicts/:id/advisory", requestAdvisory)
	v1.Get("/audit/voxel/:h3Cell/:altBin/:timeBin", getAuditTrail)

	app.Get("/healthz", func(c fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })
	app.Get("/ws/live-updates", websocket.New(func(c *websocket.Conn) {
		hub.Register(c)
		defer hub.Unregister(c)
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break // client disconnected
			}
		}
	}))
}

func submitFlightPlan(c fiber.Ctx) error {
	var req FlightPlanRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	result, err := standardsClient.SubmitFlightPlan(c.Context(), toProto(req))
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, "standards-svc unavailable")
	}
	return c.Status(fiber.StatusAccepted).JSON(fromProto(result))
}