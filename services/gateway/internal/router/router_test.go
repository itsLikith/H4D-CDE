// Copyright 2026 Likith Saragadam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/itsLikith/h4d-cde/services/gateway/internal/router"
	"github.com/itsLikith/h4d-cde/services/gateway/internal/ws"
	"github.com/stretchr/testify/assert"
)

func TestHealthzEndpoint(t *testing.T) {
	app := fiber.New()
	hub := ws.NewHub()
	router.Register(app, hub, router.Clients{})

	req := httptest.NewRequest("GET", "/healthz", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
