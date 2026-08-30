package adaptive

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

// ForecastFor reads the latest occupancy forecast demand-forecaster-svc wrote directly to Redis.
func ForecastFor(ctx context.Context, client *redis.Client, h3Cell string) (int, error) {
	raw, err := client.Get(ctx, "forecast:"+h3Cell).Result()
	if err == redis.Nil {
		return 0, nil // no forecast yet -- treat as low density, use base resolution
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
	return steps[0], nil // nearest-horizon prediction
}