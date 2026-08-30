# Copyright 2026 Likith Saragadam
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: all proto train-models test benchmark build up down clean

SHELL := /bin/bash
PYTHON := services/trajectory-predictor-svc/.venv/bin/python3

all: proto build test

## Compile Protocol Buffer definitions for Go and Python
proto:
	@echo "[*] Generating Protobuf and gRPC stubs with buf..."
	cd proto && buf generate

## Train all 3 AI models (Trajectory Predictor, Risk Scorer, Demand Forecaster)
train-models:
	@echo "[*] Training AI Augmentation models..."
	export PYTHONPATH=$$(pwd):$$(pwd)/services/trajectory-predictor-svc:$$(pwd)/services/risk-scorer-svc:$$(pwd)/services/demand-forecaster-svc:$$(pwd)/tools/synthetic-data-gen && \
	$(PYTHON) -m services.trajectory-predictor-svc.app.train && \
	$(PYTHON) -m services.risk-scorer-svc.app.train && \
	$(PYTHON) -m services.demand-forecaster-svc.app.train

## Run unit tests across all Go microservices and Python modules
test:
	@echo "[*] Running Go test suites..."
	cd services/voxel-engine && go test -v ./...
	cd services/standards-svc && go test -v ./...
	cd services/audit-svc && go test -v ./...
	cd services/gateway && go test -v ./...

## Run empirical benchmark reproducing Table I, II, and III
benchmark:
	@echo "[*] Running empirical benchmark reproducing Sahadevan et al. (ICSPIS 2025)..."
	go run ./tests/benchmark/run_benchmark.go

## Build all Go service binaries
build:
	@echo "[*] Building Go microservice binaries..."
	cd services/voxel-engine && CGO_ENABLED=1 go build -o ../../bin/voxel-engine ./cmd/server/main.go
	cd services/standards-svc && CGO_ENABLED=0 go build -o ../../bin/standards-svc ./cmd/server/main.go
	cd services/audit-svc && CGO_ENABLED=0 go build -o ../../bin/audit-svc ./cmd/server/main.go
	cd services/gateway && CGO_ENABLED=0 go build -o ../../bin/gateway ./cmd/server/main.go

## One-click start all services and infrastructure in Docker
up:
	@echo "[*] Launching H4D-CDE microservices and infrastructure..."
	docker compose up --build -d

## Tear down all containers and volumes
down:
	@echo "[*] Stopping H4D-CDE..."
	docker compose down -v

clean:
	rm -rf bin/ models/*.joblib models/*.pt
