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

// Package occupancy implements the shared 4D occupancy map O[voxel_key] = {e1, e2, ...} from Equation (7).
package occupancy

import (
	"context"
	"sync"
	"time"

	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/temporal"
	"github.com/redis/go-redis/v9"
)

// Map defines the contract for storing and querying entity presence in 4D voxel cells.
type Map interface {
	Add(ctx context.Context, key temporal.VoxelKey, entityID string) error
	Occupants(ctx context.Context, key temporal.VoxelKey) ([]string, error)
	Reset(ctx context.Context) error
}

// RedisMap is the distributed, production-grade occupancy backend backed by Redis Sets.
// Provides O(1) amortized insertion (SADD) and lookup (SMEMBERS) across all horizontally scaled replicas.
type RedisMap struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisMap constructs a Redis-backed occupancy store with auto-expiring keys.
func NewRedisMap(client *redis.Client, ttl time.Duration) *RedisMap {
	if ttl <= 0 {
		ttl = 120 * time.Second
	}
	return &RedisMap{
		client: client,
		ttl:    ttl,
	}
}

func redisKey(k temporal.VoxelKey) string {
	return "voxel:" + k.String()
}

func (m *RedisMap) Add(ctx context.Context, key temporal.VoxelKey, entityID string) error {
	rk := redisKey(key)
	pipe := m.client.Pipeline()
	pipe.SAdd(ctx, rk, entityID)
	pipe.Expire(ctx, rk, m.ttl)
	pipe.SAdd(ctx, "active_voxel_cells", key.H3Cell.String())
	_, err := pipe.Exec(ctx)
	return err
}

func (m *RedisMap) Occupants(ctx context.Context, key temporal.VoxelKey) ([]string, error) {
	return m.client.SMembers(ctx, redisKey(key)).Result()
}

func (m *RedisMap) Reset(ctx context.Context) error {
	return m.client.FlushDB(ctx).Err()
}

// InMemoryMap provides an in-memory concurrent map implementation ideal for unit testing and local benchmarks.
type InMemoryMap struct {
	mu   sync.RWMutex
	data map[string]map[string]struct{}
}

// NewInMemoryMap instantiates a thread-safe in-memory occupancy map.
func NewInMemoryMap() *InMemoryMap {
	return &InMemoryMap{
		data: make(map[string]map[string]struct{}),
	}
}

func (m *InMemoryMap) Add(_ context.Context, key temporal.VoxelKey, entityID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	k := key.String()
	if m.data[k] == nil {
		m.data[k] = make(map[string]struct{})
	}
	m.data[k][entityID] = struct{}{}
	return nil
}

func (m *InMemoryMap) Occupants(_ context.Context, key temporal.VoxelKey) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	k := key.String()
	set, exists := m.data[k]
	if !exists {
		return []string{}, nil
	}

	out := make([]string, 0, len(set))
	for e := range set {
		out = append(out, e)
	}
	return out, nil
}

func (m *InMemoryMap) Reset(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]map[string]struct{})
	return nil
}
