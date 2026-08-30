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
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	standardsv1 "github.com/itsLikith/h4d-cde/gen/standards"
	voxelizerv1 "github.com/itsLikith/h4d-cde/gen/voxelizer"

	"github.com/itsLikith/h4d-cde/services/standards-svc/internal/grpcserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("[*] Starting H4D-CDE Standards & SCD Service...")

	port := 50050
	if pStr := os.Getenv("STANDARDS_SVC_GRPC_PORT"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil {
			port = p
		}
	} else if pStr := os.Getenv("PORT"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil {
			port = p
		}
	}

	voxelEngineHost := os.Getenv("VOXEL_ENGINE_HOST")
	if voxelEngineHost == "" {
		voxelEngineHost = "localhost"
	}
	voxelEnginePort := os.Getenv("VOXEL_ENGINE_GRPC_PORT")
	if voxelEnginePort == "" {
		voxelEnginePort = "50051"
	}
	voxelEngineTarget := voxelEngineHost + ":" + voxelEnginePort

	var veClient voxelizerv1.VoxelEngineServiceClient
	if conn, err := grpc.NewClient(voxelEngineTarget, grpc.WithTransportCredentials(insecure.NewCredentials())); err == nil {
		veClient = voxelizerv1.NewVoxelEngineServiceClient(conn)
		log.Printf("[+] Connected to Voxel Engine at %s", voxelEngineTarget)
	} else {
		log.Printf("[!] Warning: Could not connect to Voxel Engine at %s: %v", voxelEngineTarget, err)
	}

	serverImpl := grpcserver.NewServer(veClient)
	grpcServer := grpc.NewServer()
	standardsv1.RegisterStandardsServiceServer(grpcServer, serverImpl)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to bind port %d: %v", port, err)
	}

	go func() {
		log.Printf("[*] Standards Service listening on gRPC port :%d", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("[*] Shutting down Standards Service...")
	grpcServer.GracefulStop()
	log.Println("[+] Standards Service shutdown complete.")
}
