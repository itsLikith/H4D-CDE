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

package storage

import (
	"context"
	"sync"

	"github.com/itsLikith/h4d-cde/services/audit-svc/internal/chain"
)

type Store interface {
	Save(ctx context.Context, entry chain.Entry) error
	GetEntries(ctx context.Context, startTimeMs, endTimeMs int64) ([]chain.Entry, error)
}

type MemoryStore struct {
	mu      sync.RWMutex
	entries []chain.Entry
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		entries: make([]chain.Entry, 0),
	}
}

func (s *MemoryStore) Save(_ context.Context, entry chain.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, entry)
	return nil
}

func (s *MemoryStore) GetEntries(_ context.Context, startTimeMs, endTimeMs int64) ([]chain.Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []chain.Entry
	for _, e := range s.entries {
		if (startTimeMs <= 0 || e.TimestampUnixMs >= startTimeMs) &&
			(endTimeMs <= 0 || e.TimestampUnixMs <= endTimeMs) {
			result = append(result, e)
		}
	}
	return result, nil
}
