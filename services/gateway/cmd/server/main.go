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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v3"
	auditv1 "github.com/itsLikith/h4d-cde/gen/audit"
	standardsv1 "github.com/itsLikith/h4d-cde/gen/standards"
	voxelizerv1 "github.com/itsLikith/h4d-cde/gen/voxelizer"

	"github.com/itsLikith/h4d-cde/services/gateway/internal/router"
	"github.com/itsLikith/h4d-cde/services/gateway/internal/ws"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("[*] Starting H4D-CDE API Gateway & WebSocket BFF...")

	port := 8080
	if pStr := os.Getenv("GATEWAY_PORT"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil {
			port = p
		}
	} else if pStr := os.Getenv("PORT"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil {
			port = p
		}
	}

	standardsHost := os.Getenv("STANDARDS_SVC_HOST")
	if standardsHost == "" {
		standardsHost = "localhost"
	}
	standardsPort := os.Getenv("STANDARDS_SVC_GRPC_PORT")
	if standardsPort == "" {
		standardsPort = "50050"
	}
	standardsTarget := standardsHost + ":" + standardsPort

	voxelEngineHost := os.Getenv("VOXEL_ENGINE_HOST")
	if voxelEngineHost == "" {
		voxelEngineHost = "localhost"
	}
	voxelEnginePort := os.Getenv("VOXEL_ENGINE_GRPC_PORT")
	if voxelEnginePort == "" {
		voxelEnginePort = "50051"
	}
	voxelEngineTarget := voxelEngineHost + ":" + voxelEnginePort

	auditHost := os.Getenv("AUDIT_SVC_HOST")
	if auditHost == "" {
		auditHost = "localhost"
	}
	auditPort := os.Getenv("AUDIT_SVC_GRPC_PORT")
	if auditPort == "" {
		auditPort = "50054"
	}
	auditTarget := auditHost + ":" + auditPort

	var clients router.Clients
	if conn, err := grpc.NewClient(standardsTarget, grpc.WithTransportCredentials(insecure.NewCredentials())); err == nil {
		clients.StandardsClient = standardsv1.NewStandardsServiceClient(conn)
		log.Printf("[+] Standards service client targeted at %s", standardsTarget)
	}

	if conn, err := grpc.NewClient(voxelEngineTarget, grpc.WithTransportCredentials(insecure.NewCredentials())); err == nil {
		clients.VoxelEngineClient = voxelizerv1.NewVoxelEngineServiceClient(conn)
		log.Printf("[+] Voxel Engine client targeted at %s", voxelEngineTarget)
	}

	if conn, err := grpc.NewClient(auditTarget, grpc.WithTransportCredentials(insecure.NewCredentials())); err == nil {
		clients.AuditClient = auditv1.NewAuditServiceClient(conn)
		log.Printf("[+] Audit service client targeted at %s", auditTarget)
	}

	app := fiber.New(fiber.Config{
		AppName: "H4D-CDE API Gateway v1.0",
	})
	hub := ws.NewHub()

	router.Register(app, hub, clients)

	kafkaBrokersStr := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersStr == "" {
		kafkaBrokersStr = "localhost:9092"
	}
	conflictsTopic := os.Getenv("KAFKA_TOPIC_CONFLICTS_DETECTED")
	if conflictsTopic == "" {
		conflictsTopic = "conflicts.detected"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := ws.BridgeKafkaToClients(ctx, []string{kafkaBrokersStr}, conflictsTopic, hub); err != nil && err != context.Canceled {
			log.Printf("[!] Gateway WebSocket Kafka bridge warning: %v", err)
		}
	}()

	go func() {
		addr := fmt.Sprintf(":%d", port)
		log.Printf("[*] API Gateway listening on HTTP port %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Fiber server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("[*] Shutting down Gateway gracefully...")
	cancel()
	_ = app.Shutdown()
	log.Println("[+] Gateway shutdown complete.")
}