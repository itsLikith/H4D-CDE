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
# standards-svc — ASTM F3548-21 / ICAO separation validation (Go / gRPC)
# ---------------------------------------------------------------------------

FROM golang:1.25.0 AS builder

ENV GOWORK=off
ENV GOTOOLCHAIN=auto
ENV CGO_ENABLED=0
ENV GOOS=linux

COPY gen/ /gen/

WORKDIR /build

COPY services/standards-svc/go.mod services/standards-svc/go.sum ./
RUN go mod download

COPY services/standards-svc/ ./
RUN go mod tidy && go build -o /app/standards-svc ./cmd/server/main.go

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/standards-svc /app/standards-svc

EXPOSE 50050
ENTRYPOINT ["/app/standards-svc"]
