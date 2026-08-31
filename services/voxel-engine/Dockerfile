# syntax=docker/dockerfile:1
# Copyright 2026 H4D-CDE Authors
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

# ---------------------------------------------------------------------------
# voxel-engine — Uber H3 hexagonal 4D voxelization engine (Go + CGO)
#
# CGO is required because h3-go wraps the Uber H3 C library.
# The builder image ships GCC. The final stage is debian-slim (not distroless)
# because the CGO binary is dynamically linked against libc.
# ---------------------------------------------------------------------------

FROM golang:1.25.0 AS builder

ENV GOWORK=off
ENV GOTOOLCHAIN=auto
ENV CGO_ENABLED=1
ENV GOOS=linux

# Install C build tooling for the h3-go CGO bindings
RUN apt-get update && apt-get install -y --no-install-recommends \
    cmake \
    gcc \
    && rm -rf /var/lib/apt/lists/*

COPY gen/ /gen/

WORKDIR /build

COPY services/voxel-engine/go.mod services/voxel-engine/go.sum ./
RUN go mod download

COPY services/voxel-engine/ ./
RUN go mod tidy && go build -o /app/voxel-engine ./cmd/server/main.go

# ---------------------------------------------------------------------------
# Final image — debian-slim carries libc for the CGO binary
# ---------------------------------------------------------------------------
FROM debian:trixie-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/voxel-engine /app/voxel-engine

EXPOSE 50051
ENTRYPOINT ["/app/voxel-engine"]
