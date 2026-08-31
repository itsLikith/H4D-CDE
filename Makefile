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
#   make setup          Bootstrap Python virtual environments (run once after clone)
#   make all            proto → build → test  (full CI pipeline)
#   make proto          Regenerate gRPC / Protobuf stubs via buf
#   make build          Compile all Go service binaries to ./bin/
#   make test           Run all Go and Python unit tests
#   make train-models   Train Trajectory Predictor, Risk Scorer, Demand Forecaster
#   make benchmark      Empirical benchmark reproducing paper Tables I–III
#   make up             Build Docker images and start the full 11-container stack
#   make down           Stop and remove all containers and volumes
#   make clean          Remove compiled binaries and model checkpoints
# ==============================================================================

.PHONY: all setup proto train-models test benchmark build up down clean

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
# setup — Bootstrap Python virtual environments (run once after git clone)
# ---------------------------------------------------------------------------
## Create and populate Python venvs for all three ML services
setup:
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

	@echo "[+] Python environments ready."

# ---------------------------------------------------------------------------
# proto — Regenerate gRPC / Protobuf stubs for Go and Python
# ---------------------------------------------------------------------------
## Compile Protocol Buffer definitions for Go and Python via buf
proto:
	@echo "[*] Generating Protobuf and gRPC stubs with buf..."
	cd proto && buf generate
	@echo "[+] Proto stubs generated."

# ---------------------------------------------------------------------------
# build — Compile all Go service binaries
# ---------------------------------------------------------------------------
## Compile all four Go microservices to ./bin/
build:
	@echo "[*] Building Go microservice binaries..."
	@mkdir -p bin
	cd services/voxel-engine   && CGO_ENABLED=1 go build -o ../../bin/voxel-engine   ./cmd/server/main.go
	cd services/standards-svc  && CGO_ENABLED=0 go build -o ../../bin/standards-svc  ./cmd/server/main.go
	cd services/audit-svc      && CGO_ENABLED=0 go build -o ../../bin/audit-svc      ./cmd/server/main.go
	cd services/gateway        && CGO_ENABLED=0 go build -o ../../bin/gateway        ./cmd/server/main.go
	@echo "[+] Binaries written to ./bin/"

# ---------------------------------------------------------------------------
# test — Run all unit tests (Go + Python)
# ---------------------------------------------------------------------------
## Run Go and Python test suites
test: _check_venvs
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
# up / down — Docker Compose lifecycle
# ---------------------------------------------------------------------------
## Build Docker images and start the full 11-container stack
up:
	@echo "[*] Launching H4D-CDE microservices and infrastructure..."
	docker compose up --build -d
	@echo "[+] Stack is up. Gateway: http://localhost:8080  Prometheus: http://localhost:9090"

## Stop and remove all containers and volumes
down:
	@echo "[*] Stopping H4D-CDE..."
	docker compose down -v
	@echo "[+] All containers and volumes removed."

# ---------------------------------------------------------------------------
# clean — Remove build artifacts
# ---------------------------------------------------------------------------
## Remove compiled binaries and trained model checkpoints
clean:
	@rm -rf bin/
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
