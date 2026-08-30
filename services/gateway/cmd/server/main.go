// services/gateway/main.go
package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"hive/gateway/internal/router"
	"hive/gateway/internal/ws"
)

func main() {
	app := fiber.New()
	hub := ws.NewHub()

	router.Register(app, hub)

	go func() {
		if err := ws.BridgeKafkaToClients(context.Background(), []string{"kafka:9092"}, hub); err != nil {
			log.Printf("kafka bridge stopped: %v", err)
		}
	}()

	log.Fatal(app.Listen(":8080"))
}