// services/voxel-engine/internal/occupancy/occupancy.go
package occupancy

import (
	"context"
	"hive/voxel-engine/internal/temporal"
	"time"

	"github.com/redis/go-redis/v9"
)

type Map interface {
	Add(ctx context.Context, key temporal.VoxelKey, entityID string) error
	Occupants(ctx context.Context, key temporal.VoxelKey) ([]string, error)
}

// RedisMap is the production backend. Every voxel key becomes a Redis SET,
// which is what makes voxel-engine horizontally scalable — every
// replica reads and writes the same shared state instead of its own memory.
type RedisMap struct {
	client *redis.Client
	ttl    time.Duration // auto-expiry: a voxel's time bin is only ever
	// relevant for the 10s window it represents
}

func NewRedisMap(client *redis.Client, ttl time.Duration) *RedisMap {
	return &RedisMap{client: client, ttl: ttl}
}

func redisKey(k temporal.VoxelKey) string { return "voxel:" + k.String() }

func (m *RedisMap) Add(ctx context.Context, key temporal.VoxelKey, entityID string) error {
	rk := redisKey(key)
	if err := m.client.SAdd(ctx, rk, entityID).Err(); err != nil {
		return err
	}
	return m.client.Expire(ctx, rk, m.ttl).Err()
}

func (m *RedisMap) Occupants(ctx context.Context, key temporal.VoxelKey) ([]string, error) {
	return m.client.SMembers(ctx, redisKey(key)).Result()
}

// InMemoryMap is the unit-test backend — not shared across processes.
type InMemoryMap struct {
	data map[string]map[string]struct{}
}

func NewInMemoryMap() *InMemoryMap { return &InMemoryMap{data: map[string]map[string]struct{}{}} }

func (m *InMemoryMap) Add(_ context.Context, key temporal.VoxelKey, entityID string) error {
	k := key.String()
	if m.data[k] == nil {
		m.data[k] = map[string]struct{}{}
	}
	m.data[k][entityID] = struct{}{}
	return nil
}

func (m *InMemoryMap) Occupants(_ context.Context, key temporal.VoxelKey) ([]string, error) {
	out := make([]string, 0, len(m.data[key.String()]))
	for e := range m.data[key.String()] {
		out = append(out, e)
	}
	return out, nil
}
