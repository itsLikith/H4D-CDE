<div align="center">

<h1>H4D-CDE</h1>
<h3>Hexagonal 4D Conflict-Detection Engine</h3>
<p><em>AI-Augmented Hexagonal Voxelization for Scalable Conflict Detection in Urban Air Mobility</em></p>

<p>
  <a href="https://doi.org/10.1109/ICSPIS67605.2025.11318403">
    <img src="https://img.shields.io/badge/Research%20Paper-ICSPIS%202025-0A66C2?style=flat-square" alt="Research Paper">
  </a>
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue?style=flat-square" alt="Apache 2.0">
  <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Python-3.12-3776AB?style=flat-square&logo=python" alt="Python">
  <img src="https://img.shields.io/badge/Next.js-16-black?style=flat-square&logo=next.js" alt="Next.js">
  <img src="https://img.shields.io/badge/Tailwind-CSS%20v4-38B2AC?style=flat-square&logo=tailwind-css" alt="Tailwind CSS">
  <img src="https://img.shields.io/badge/shadcn%2Fui-Base%20UI-000000?style=flat-square" alt="shadcn/ui">
  <img src="https://img.shields.io/badge/Docker-Compose-2496ED?style=flat-square&logo=docker" alt="Docker">
  <img src="https://img.shields.io/badge/gRPC-protobuf-5FB8E6?style=flat-square" alt="gRPC">
</p>

</div>

---

## Overview

H4D-CDE is a production-grade, open-source implementation of the methods described in:

> **AI-Augmented Hexagonal Voxelization for Scalable Conflict Detection in Urban Air Mobility**  
> Sahadevan, Al Ali & Mahesh — ICSPIS 2025  
> DOI: [10.1109/ICSPIS67605.2025.11318403](https://doi.org/10.1109/ICSPIS67605.2025.11318403)

The system decomposes urban airspace into a hierarchical 4D grid of hexagonal voxels (Uber H3 resolution 8, ~0.74 km² cells × 100 ft altitude bins × 10 s time bins) and uses a pipeline of three AI models to proactively detect separation violations before they occur.

### Architecture

```
  Flight Plan ──► [Standards Svc]──►[Voxel Engine]──►[Trajectory Predictor]
                   ASTM F3548-21     H3 grid cells     RF / kinematic 4D path
                   validation        + Redis cache
                                          │
                                ┌─────────▼──────────┐
                                │  Demand Forecaster  │  (TCN 15-min density)
                                └─────────┬──────────┘
                                          │
                                ┌─────────▼──────────┐
                                │    Risk Scorer      │  (XGBoost P(conflict))
                                └─────────┬──────────┘
                                          │
 [Web Dashboard :3000] ◄── [Gateway :8080] ◄──── │ ──── [Audit Svc → Kafka → Postgres]
```

**Services and Infrastructure:**

| Component | Technology | Port | Role |
|---|---|---|---|
| `web` | Next.js 16 / React 19 / Tailwind / shadcn | `:3000` | Situational Awareness & Mission Control UI |
| `gateway` | Go / Fiber v3 | `:8080` | REST + WebSocket BFF, JWT auth |
| `standards-svc` | Go / gRPC | `:50050` | ASTM F3548-21 / ICAO separation validation |
| `voxel-engine` | Go / Uber H3 CGO | `:50051` | 4D hexagonal voxelization, Redis occupancy store |
| `trajectory-predictor-svc` | Python / scikit-learn | `:50052` | Gradient Boosting 4D path inference |
| `demand-forecaster-svc` | Python / PyTorch | `:50053` | Dilated TCN 15-min airspace density forecast |
| `audit-svc` | Go / gRPC | `:50054` | Cryptographic SHA-256 hash-chain audit ledger |
| `risk-scorer-svc` | Python / XGBoost | `:50055` | Pairwise conflict probability scoring |
| `postgres` | PostgreSQL 16 | `:5432` | Durable audit log & flight plan persistence |
| `redis` | Redis 7.2 | `:6379` | Spatial occupancy store & forecast cache |
| `kafka` | Redpanda | `:9092` | Distributed audit event & conflict stream |
| `prometheus` | Prometheus v2.50 | `:9090` | System metrics scraper & health monitoring |

---

## Prerequisites

| Tool | Minimum version | Install |
|---|---|---|
| **Docker Desktop** | 4.x | https://docs.docker.com/get-docker/ |
| **Docker Compose** | v2.x | bundled with Docker Desktop |
| **Go** | 1.25 | https://go.dev/dl/ |
| **Python** | 3.12 | https://www.python.org/downloads/ |
| **Node.js** | 20+ / 22+ | https://nodejs.org/ |
| **buf** | 1.x | `brew install bufbuild/buf/buf` |
| **Make** | any | pre-installed on macOS / Linux |

---

## Quick Start (Docker)

### 1. Clone & Configure

```bash
git clone https://github.com/itsLikith/H4D-CDE.git
cd H4D-CDE
cp .env.example .env
```

### 2. Launch the Full Stack

```bash
make up
```

This builds all images and starts the **12-container stack**:
- **Web Dashboard**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Prometheus**: http://localhost:9090

### 3. Verify Health

```bash
docker compose ps
curl http://localhost:8080/healthz
# {"engine":"H4D-CDE","service":"gateway","status":"ok"}
```

---

## Development Setup

### Bootstrap All Environments

```bash
make setup
```
This automatically:
1. Installs Node.js dependencies in `web/` (`npm ci`).
2. Creates and populates isolated Python virtual environments for all three ML services (`trajectory-predictor-svc`, `risk-scorer-svc`, `demand-forecaster-svc`).

### Build Binaries & Web Dashboard

```bash
make build
```

### Run Tests & Verification

```bash
make test
```
Runs:
- Go unit test suites (`voxel-engine`, `standards-svc`, `audit-svc`, `gateway`).
- Python pytest suites across all 3 ML services.
- Web dashboard ESLint suite.

### Run Empirical Benchmark

Reproduces Tables I, II, and III from Sahadevan et al. (ICSPIS 2025):

```bash
make benchmark
```

### Start Web Dashboard Locally

```bash
make web-dev
# Launches Next.js on http://localhost:3000
```

---

## Makefile Reference

| Target | Description |
|---|---|
| `make all` | `proto` → `build` → `test` (full CI pipeline) |
| `make setup` | Bootstrap Python venvs & Node.js dependencies |
| `make proto` | Regenerate gRPC / Protobuf stubs via `buf` |
| `make build` | Compile Go binaries to `./bin/` & build Next.js dashboard |
| `make test` | Run Go, Python, and Web test suites |
| `make train-models` | Train Trajectory Predictor, Risk Scorer, and Demand Forecaster |
| `make benchmark` | Run empirical benchmark reproducing paper Tables I–III |
| `make web-dev` | Start local Next.js dashboard on `:3000` |
| `make web-lint` | Lint web dashboard with ESLint |
| `make web-format` | Format web code with Prettier |
| `make up` | Build Docker images & launch the full 12-container stack |
| `make down` | Stop and remove all containers and volumes |
| `make clean` | Clean compiled binaries, web build artifacts, and model weights |

---

## License & Citation

Apache License 2.0 — see [`LICENSE`](LICENSE).

```bibtex
@inproceedings{sahadevan2025ai,
  title     = {AI-Augmented Hexagonal Voxelization for Scalable Conflict Detection in Urban Air Mobility},
  author    = {Sahadevan, Deepudev and Al Ali, Hannah and Mahesh, Chilka},
  booktitle = {2025 8th International Conference on Signal Processing and Information Security (ICSPIS)},
  year      = {2025},
  doi       = {10.1109/ICSPIS67605.2025.11318403}
}
```
