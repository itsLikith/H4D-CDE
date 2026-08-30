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
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	auditv1 "github.com/itsLikith/h4d-cde/gen/audit"

	"github.com/itsLikith/h4d-cde/services/audit-svc/internal/chain"
	"github.com/itsLikith/h4d-cde/services/audit-svc/internal/consumer"
	"github.com/itsLikith/h4d-cde/services/audit-svc/internal/grpcserver"
	"github.com/itsLikith/h4d-cde/services/audit-svc/internal/storage"
	"google.golang.org/grpc"
)

func main() {
	log.Println("[*] Starting H4D-CDE Tamper-Evident Audit Service...")

	port := 50054
	if pStr := os.Getenv("AUDIT_SVC_GRPC_PORT"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil {
			port = p
		}
	} else if pStr := os.Getenv("PORT"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil {
			port = p
		}
	}

	auditChain := chain.New()
	store := storage.NewMemoryStore()

	kafkaBrokersStr := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersStr == "" {
		kafkaBrokersStr = "localhost:9092"
	}
	kafkaTopic := os.Getenv("KAFKA_TOPIC_AUDIT_EVENTS")
	if kafkaTopic == "" {
		kafkaTopic = "audit.events"
	}
	kafkaGroup := os.Getenv("KAFKA_CONSUMER_GROUP_AUDIT")
	if kafkaGroup == "" {
		kafkaGroup = "audit-svc-group"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Spawn Kafka Consumer loop in background
	go func() {
		err := consumer.Run(
			ctx,
			[]string{kafkaBrokersStr},
			kafkaTopic,
			kafkaGroup,
			auditChain,
			func(e chain.Entry) error {
				return store.Save(ctx, e)
			},
		)
		if err != nil && err != context.Canceled {
			log.Printf("[!] Audit Kafka consumer stopped: %v", err)
		}
	}()

	serverImpl := grpcserver.NewServer(auditChain, store)
	grpcServer := grpc.NewServer()
	auditv1.RegisterAuditServiceServer(grpcServer, serverImpl)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to bind port %d: %v", port, err)
	}

	go func() {
		log.Printf("[*] Audit Service listening on gRPC port :%d", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("[*] Shutting down Audit Service...")
	cancel()
	grpcServer.GracefulStop()
	log.Println("[+] Audit Service shutdown complete.")
}
