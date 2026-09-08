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

# ==============================================================================
# H4D-CDE — Hexagonal 4D Conflict-Detection Engine
# ==============================================================================
# Usage:
#   make setup          Bootstrap Python venvs & Node dependencies
#   make all            proto → build → test  (full CI pipeline)
#   make proto          Regenerate gRPC / Protobuf stubs via buf
#   make build          Compile all Go binaries & build Next.js web dashboard
#   make test           Run all Go, Python, and Frontend test suites
#   make train-models   Train Trajectory Predictor, Risk Scorer, Demand Forecaster
#   make benchmark      Empirical benchmark reproducing paper Tables I–III
#   make web-dev        Start Next.js web dashboard in local development mode
#   make up             Build Docker images and start the full 12-container stack
#   make down           Stop and remove all containers and volumes
#   make clean          Remove compiled binaries, web build, and model checkpoints
# ==============================================================================

.PHONY: all setup proto train-models test benchmark build up down clean \
        web-setup web-build web-lint web-format web-dev

SHELL := /bin/bash

# ---------------------------------------------------------------------------
# Python interpreter paths — each service carries its own isolated venv so
# dependency versions never conflict across ML frameworks.
# ---------------------------------------------------------------------------
PYTHON_TRAJ   := services/trajectory-predictor-svc/.venv/bin/python3
PYTHON_RISK   := services/risk-scorer-svc/.venv/bin/python3
PYTHON_DEMAND := services/demand-forecaster-svc/.venv/bin/python3

# ---------------------------------------------------------------------------
# Top-level targets
# ---------------------------------------------------------------------------

## Run the full CI pipeline: proto → build → test
all: proto build test

# ---------------------------------------------------------------------------
# setup — Bootstrap environments (run once after git clone)
# ---------------------------------------------------------------------------
## Create and populate Python venvs & install Node.js web dashboard dependencies
setup: web-setup
	@echo "[*] Bootstrapping Python virtual environments..."
	@if [ ! -f services/trajectory-predictor-svc/.venv/bin/python3 ]; then \
		echo "[*] Creating venv: trajectory-predictor-svc"; \
		python3 -m venv services/trajectory-predictor-svc/.venv; \
	fi
	services/trajectory-predictor-svc/.venv/bin/pip install --quiet --upgrade pip
	services/trajectory-predictor-svc/.venv/bin/pip install --quiet -r services/trajectory-predictor-svc/requirements.txt

	@if [ ! -f services/risk-scorer-svc/.venv/bin/python3 ]; then \
		echo "[*] Creating venv: risk-scorer-svc"; \
		python3 -m venv services/risk-scorer-svc/.venv; \
	fi
	services/risk-scorer-svc/.venv/bin/pip install --quiet --upgrade pip
	services/risk-scorer-svc/.venv/bin/pip install --quiet -r services/risk-scorer-svc/requirements.txt

	@if [ ! -f services/demand-forecaster-svc/.venv/bin/python3 ]; then \
		echo "[*] Creating venv: demand-forecaster-svc"; \
		python3 -m venv services/demand-forecaster-svc/.venv; \
	fi
	services/demand-forecaster-svc/.venv/bin/pip install --quiet --upgrade pip
	services/demand-forecaster-svc/.venv/bin/pip install --quiet -r services/demand-forecaster-svc/requirements.txt

	@echo "[+] Python and Web environments ready."

# ---------------------------------------------------------------------------
# proto — Regenerate gRPC / Protobuf stubs for Go and Python
# ---------------------------------------------------------------------------
## Compile Protocol Buffer definitions for Go and Python via buf
proto:
	@echo "[*] Generating Protobuf and gRPC stubs with buf..."
	cd proto && buf generate
	@echo "[+] Proto stubs generated."

# ---------------------------------------------------------------------------
# build — Compile all Go service binaries & Web dashboard
# ---------------------------------------------------------------------------
## Compile all four Go microservices to ./bin/ and build Next.js dashboard
build: web-build
	@echo "[*] Building Go microservice binaries..."
	@mkdir -p bin
	cd services/voxel-engine   && CGO_ENABLED=1 go build -o ../../bin/voxel-engine   ./cmd/server/main.go
	cd services/standards-svc  && CGO_ENABLED=0 go build -o ../../bin/standards-svc  ./cmd/server/main.go
	cd services/audit-svc      && CGO_ENABLED=0 go build -o ../../bin/audit-svc      ./cmd/server/main.go
	cd services/gateway        && CGO_ENABLED=0 go build -o ../../bin/gateway        ./cmd/server/main.go
	@echo "[+] Binaries written to ./bin/"

# ---------------------------------------------------------------------------
# test — Run all unit tests (Go + Python + Frontend Lint)
# ---------------------------------------------------------------------------
## Run Go, Python, and Web test/lint suites
test: _check_venvs web-lint
	@echo "[*] Running Go test suites..."
	cd services/voxel-engine  && go test -v ./...
	cd services/standards-svc && go test -v ./...
	cd services/audit-svc     && go test -v ./...
	cd services/gateway       && go test -v ./...

	@echo "[*] Running Python test suites..."
	cd services/trajectory-predictor-svc && $(abspath $(PYTHON_TRAJ)) -m pytest tests/ -v
	cd services/risk-scorer-svc          && $(abspath $(PYTHON_RISK)) -m pytest tests/ -v
	cd services/demand-forecaster-svc    && $(abspath $(PYTHON_DEMAND)) -m pytest tests/ -v
	@echo "[+] All tests passed."

# ---------------------------------------------------------------------------
# train-models — Train all three AI augmentation models
# ---------------------------------------------------------------------------
## Train Trajectory Predictor (GBM), Risk Scorer (XGBoost), Demand Forecaster (TCN)
train-models: _check_venvs
	@echo "[*] Training AI Augmentation models..."
	@mkdir -p models

	@echo "[*] Training Trajectory Predictor (Gradient Boosting Regressor)..."
	cd services/trajectory-predictor-svc && $(abspath $(PYTHON_TRAJ)) -m app.train

	@echo "[*] Training Risk Scorer (XGBoost Classifier)..."
	cd services/risk-scorer-svc && $(abspath $(PYTHON_RISK)) -m app.train

	@echo "[*] Training Demand Forecaster (Dilated TCN)..."
	cd services/demand-forecaster-svc && $(abspath $(PYTHON_DEMAND)) -m app.train

	@echo "[+] Model checkpoints written to ./models/"

# ---------------------------------------------------------------------------
# benchmark — Empirical benchmark reproducing paper Tables I, II, III
# ---------------------------------------------------------------------------
## Reproduce Sahadevan et al. (ICSPIS 2025) benchmark results
benchmark:
	@echo "[*] Running empirical benchmark (Sahadevan et al. ICSPIS 2025)..."
	go run ./tests/benchmark/run_benchmark.go
	@echo "[+] Benchmark complete."

# ---------------------------------------------------------------------------
# Frontend / Web Dashboard Targets
# ---------------------------------------------------------------------------
## Install web dependencies
web-setup:
	@echo "[*] Installing web dependencies..."
	cd web && npm ci

## Build production Next.js web application
web-build:
	@echo "[*] Building Next.js web dashboard..."
	cd web && npm run build

## Lint web application
web-lint:
	@echo "[*] Linting web dashboard..."
	cd web && npm run lint

## Format web application code with Prettier
web-format:
	@echo "[*] Formatting web dashboard code with Prettier..."
	cd web && npm run format

## Start local Next.js web dashboard in dev mode on port 3000
web-dev:
	@echo "[*] Starting H4D-CDE Web Dashboard on http://localhost:3000..."
	cd web && npm run dev -- -p 3000

# ---------------------------------------------------------------------------
# up / down — Docker Compose lifecycle
# ---------------------------------------------------------------------------
## Build Docker images and start the full 12-container stack
up:
	@echo "[*] Launching H4D-CDE microservices, web dashboard, and infrastructure..."
	docker compose up --build -d
	@echo "[+] Stack is up."
	@echo "    - Web Dashboard:  http://localhost:3000"
	@echo "    - API Gateway:    http://localhost:8080"
	@echo "    - Prometheus:     http://localhost:9090"

## Stop and remove all containers and volumes
down:
	@echo "[*] Stopping H4D-CDE..."
	docker compose down -v
	@echo "[+] All containers and volumes removed."

# ---------------------------------------------------------------------------
# clean — Remove build artifacts
# ---------------------------------------------------------------------------
## Remove compiled binaries, web build artifacts, and trained model checkpoints
clean:
	@rm -rf bin/ web/.next/ web/out/
	@rm -f models/*.joblib models/*.pt
	@echo "[+] Clean complete."

# ---------------------------------------------------------------------------
# Internal helpers (not shown in `make help`)
# ---------------------------------------------------------------------------

# Guard that ensures Python venvs are populated before running Python targets.
# Surfaces a clear error message instead of a confusing ModuleNotFoundError.
_check_venvs:
	@if [ ! -f "$(PYTHON_TRAJ)" ] || ! "$(PYTHON_TRAJ)" -c "import joblib" 2>/dev/null; then \
		echo ""; \
		echo "ERROR: Python venvs are not set up or packages are missing."; \
		echo "       Run  make setup  first."; \
		echo ""; \
		exit 1; \
	fi
