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

package adaptive

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

// ForecastFor reads the latest occupancy forecast computed asynchronously by demand-forecaster-svc.
// Fetches from Redis key 'forecast:<h3_cell>' with zero RPC overhead on the critical path.
func ForecastFor(ctx context.Context, client *redis.Client, h3Cell string) (int, error) {
	if client == nil || h3Cell == "" {
		return 0, nil
	}

	raw, err := client.Get(ctx, "forecast:"+h3Cell).Result()
	if err == redis.Nil {
		return 0, nil // Low density baseline
	}
	if err != nil {
		return 0, err
	}

	var steps []int
	if err := json.Unmarshal([]byte(raw), &steps); err != nil {
		return 0, err
	}
	if len(steps) == 0 {
		return 0, nil
	}

	return steps[0], nil
}
