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
# gateway — REST + WebSocket API Gateway (Go / Fiber v3)
#
# The go.mod has:  replace github.com/itsLikith/h4d-cde/gen => ../../gen
# Inside the container WORKDIR=/build, so "../../gen" resolves to /gen.
# We therefore COPY gen/ to /gen BEFORE go mod download so the replace
# directive can be satisfied during dependency resolution.
# ---------------------------------------------------------------------------

FROM golang:1.25.0 AS builder

ENV GOWORK=off
ENV GOTOOLCHAIN=auto
ENV CGO_ENABLED=0
ENV GOOS=linux

# Copy the generated proto stubs to /gen so the replace directive resolves
COPY gen/ /gen/

WORKDIR /build

# Copy module manifest and download dependencies (cached layer)
COPY services/gateway/go.mod services/gateway/go.sum ./
RUN go mod download

# Copy full service source
COPY services/gateway/ ./

# Tidy resolves any host-workspace vs container discrepancies
RUN go mod tidy && go build -o /app/gateway ./cmd/server/main.go

# ---------------------------------------------------------------------------
# Final minimal image (no shell, no package manager, just the binary)
# ---------------------------------------------------------------------------
FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/gateway /app/gateway

EXPOSE 8080
ENTRYPOINT ["/app/gateway"]
