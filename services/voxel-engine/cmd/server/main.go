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
	"syscall"
	"time"

	riskscorerv1 "github.com/itsLikith/h4d-cde/gen/riskscorer"
	trajectorypredictorv1 "github.com/itsLikith/h4d-cde/gen/trajectorypredictor"
	voxelizerv1 "github.com/itsLikith/h4d-cde/gen/voxelizer"

	"github.com/itsLikith/h4d-cde/services/voxel-engine/config"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/adaptive"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/advisory"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/grpcserver"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/occupancy"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("[*] Starting H4D-CDE Voxel Engine Service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 1. Connect to Redis Occupancy Store
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	var occ occupancy.Map
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("[!] Redis unreachable at %s (%v). Initializing in-memory occupancy fallback for local mode.", cfg.RedisAddr, err)
		occ = occupancy.NewInMemoryMap()
	} else {
		log.Printf("[+] Connected to Redis occupancy store at %s", cfg.RedisAddr)
		occ = occupancy.NewRedisMap(rdb, cfg.RedisTTL)
	}
	cancel()

	// 2. Initialize gRPC clients to downstream ML services
	var tpClient trajectorypredictorv1.TrajectoryPredictorServiceClient
	if conn, err := grpc.NewClient(cfg.TrajectoryPredictorTarget, grpc.WithTransportCredentials(insecure.NewCredentials())); err == nil {
		tpClient = trajectorypredictorv1.NewTrajectoryPredictorServiceClient(conn)
		log.Printf("[+] Trajectory Predictor target set to %s", cfg.TrajectoryPredictorTarget)
	}

	var rsClient riskscorerv1.RiskScorerServiceClient
	if conn, err := grpc.NewClient(cfg.RiskScorerTarget, grpc.WithTransportCredentials(insecure.NewCredentials())); err == nil {
		rsClient = riskscorerv1.NewRiskScorerServiceClient(conn)
		log.Printf("[+] Risk Scorer target set to %s", cfg.RiskScorerTarget)
	}

	// 3. Initialize Adaptive Discretization Engine & Advisory Selector
	adaptEngine := adaptive.New(adaptive.Config{
		BaseResolution:   cfg.Voxelization.H3Resolution,
		FineResolution:   cfg.AdaptiveDiscretization.FineResolution,
		BaseTimeBinS:     cfg.Voxelization.TimeBinS,
		FineTimeBinS:     cfg.AdaptiveDiscretization.FineTimeBinS,
		DensityThreshold: cfg.AdaptiveDiscretization.DensityThreshold,
	})

	advWeights := advisory.Weights{
		Delay:               cfg.Advisory.Weights.Delay,
		PathDeviation:       cfg.Advisory.Weights.PathDeviation,
		AltitudeChange:      cfg.Advisory.Weights.AltitudeChange,
		ConflictProbability: cfg.Advisory.Weights.ConflictProbability,
	}

	// 4. Initialize Kafka Event Publisher
	var publisher grpcserver.EventPublisher
	if len(cfg.KafkaBrokers) > 0 {
		publisher = grpcserver.NewKafkaPublisher(cfg.KafkaBrokers, "audit.events", "conflicts.detected")
	} else {
		publisher = &grpcserver.NoOpPublisher{}
	}

	// 5. Build and run gRPC Server
	serverImpl := grpcserver.NewServer(
		occ,
		rdb,
		tpClient,
		rsClient,
		adaptEngine,
		advWeights,
		cfg.Risk.AdvisoryThreshold,
		cfg.Separation.HorizontalNM,
		cfg.Separation.VerticalFt,
		publisher,
	)

	grpcServer := grpc.NewServer()
	voxelizerv1.RegisterVoxelEngineServiceServer(grpcServer, serverImpl)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to bind port %d: %v", cfg.GRPCPort, err)
	}

	go func() {
		log.Printf("[*] Voxel Engine gRPC server listening on port :%d", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Graceful shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("[*] Shutting down Voxel Engine gracefully...")
	grpcServer.GracefulStop()
	_ = publisher.Close()
	log.Println("[+] Voxel Engine shutdown complete.")
}
