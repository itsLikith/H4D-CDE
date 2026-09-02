<div align="center">

<!-- Banner -->
<img src="assets/banner.jpg" alt="H4D-CDE Banner" height="auto" width="auto">

<br><br>

<h1>H4D-CDE</h1>
<h3>Hexagonal 4D Conflict-Detection Engine</h3>
<p><em>AI-Augmented Hexagonal Voxelization for Scalable Conflict Detection in Urban Air Mobility</em></p>

<br>

<p>
  <img src="assets/emirates.jpeg" alt="Emirates Aviation University" height="52">
  &nbsp;&nbsp;&nbsp;&nbsp;
  <img src="assets/AAI.png" alt="Airports Authority of India" height="52">
</p>
<sub>Supported by Emirates Aviation University & Airports Authority of India</sub>

<br><br>

<p>
  <a href="https://doi.org/10.1109/ICSPIS67605.2025.11318403">
    <img src="https://img.shields.io/badge/Research%20Paper-ICSPIS%202025-0A66C2?style=flat-square" alt="Research Paper">
  </a>
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue?style=flat-square" alt="Apache 2.0">
  <img src="https://img.shields.io/badge/Docker-Compose-2496ED?style=flat-square&logo=docker" alt="Docker">
</p>

</div>

---

## Overview

**H4D-CDE** is a production-grade, open-source system for scalable conflict detection in Urban Air Mobility.  

It transforms continuous airspace into a hierarchical 4D hexagonal lattice and combines it with AI-driven prediction, risk scoring, and advisory generation. The result is a framework that remains computationally efficient and operationally reliable even under high traffic density.

The methods implemented in this repository are described in:

> **AI-Augmented Hexagonal Voxelization for Scalable Conflict Detection in Urban Air Mobility**  
> Sahadevan, Al Ali & Mahesh — ICSPIS 2025  
> DOI: [10.1109/ICSPIS67605.2025.11318403](https://doi.org/10.1109/ICSPIS67605.2025.11318403)

---

## Prerequisites

| Tool               | Minimum Version | Install                             |
|--------------------|-----------------|-------------------------------------|
| **Docker Desktop** | 4.x             | https://docs.docker.com/get-docker/ |
| **Docker Compose** | v2.x            | Bundled with Docker Desktop         |
| **Go**             | 1.25+           | https://go.dev/dl/                  |
| **Python**         | 3.12            | https://www.python.org/downloads/   |
| **Node.js**        | 20+ / 22+       | https://nodejs.org/                 |
| **Make**           | any             | Pre-installed on macOS / Linux      |

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

This builds all images and starts the complete multi-container environment:

- **Web Dashboard** → http://localhost:3000  
- **API Gateway** → http://localhost:8080

### 3. Verify Health

```bash
docker compose ps
curl http://localhost:8080/healthz
```

Expected response:

```json
{"engine":"H4D-CDE","service":"gateway","status":"ok"}
```

---

## Development Setup

### Bootstrap Environments

```bash
make setup
```

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Start Web Dashboard Locally

```bash
make web-dev
```

---

## Makefile Reference

| Target          | Description                                                      |
|-----------------|------------------------------------------------------------------|
| `make setup`    | Bootstrap Python virtual environments & Node.js dependencies     |
| `make build`    | Compile binaries and build the web dashboard                     |
| `make test`     | Run Go, Python, and web test suites                              |
| `make up`       | Build images and launch the full Docker stack                    |
| `make down`     | Stop and remove all containers and volumes                       |
| `make clean`    | Remove build artifacts and temporary files                       |
| `make web-dev`  | Start the Next.js dashboard in development mode                  |

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

---

<div align="center">
<sub>© 2026 — Likith Saragadam</sub>
</div>
