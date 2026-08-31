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
# risk-scorer-svc
# Python 3.12-slim image running the XGBoost conflict risk scoring gRPC service.
# ---------------------------------------------------------------------------

FROM python:3.12-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

COPY services/risk-scorer-svc/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY services/risk-scorer-svc/gen/ ./gen/
COPY services/risk-scorer-svc/app/ ./app/

RUN mkdir -p ./models

ENV PYTHONPATH=/app
CMD ["python", "-m", "app.server"]
