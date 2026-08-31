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
# trajectory-predictor-svc
# Python 3.12-slim image running the 4D trajectory inference gRPC service.
# Build context is the repo root so COPY can reach gen/ and app/.
# ---------------------------------------------------------------------------

FROM python:3.12-slim

WORKDIR /app

# System dependencies (curl for healthchecks, ca-certs for TLS)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies first (cached layer if requirements unchanged)
COPY services/trajectory-predictor-svc/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy generated protobuf stubs
COPY services/trajectory-predictor-svc/gen/ ./gen/

# Copy application source
COPY services/trajectory-predictor-svc/app/ ./app/

# Ensure models directory exists (checkpoints are mounted or generated at runtime)
RUN mkdir -p ./models

ENV PYTHONPATH=/app
CMD ["python", "-m", "app.server"]
