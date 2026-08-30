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

// Package chain implements a tamper-evident, SHA-256 hash-chained cryptographic ledger.
// Every conflict detection, voxel write, and advisory generation event is cryptographically linked
// to provide non-repudiable audit trails for aviation regulatory compliance (ASTM F3548-21 & ICAO).
package chain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"sync"
	"time"
)

// Entry represents a single immutable block in the audit chain.
type Entry struct {
	TimestampUnixMs int64           `json:"ts"`
	EventType       string          `json:"event_type,omitempty"`
	Event           json.RawMessage `json:"event"`
	PrevHash        string          `json:"prev_hash"`
	Hash            string          `json:"hash,omitempty"`
}

// Chain maintains the running head of the cryptographic ledger in memory.
type Chain struct {
	mu       sync.RWMutex
	prevHash string
	entries  []Entry
}

// New constructs an audit chain initialized with a standard 64-zero genesis hash.
func New() *Chain {
	return &Chain{
		prevHash: strings.Repeat("0", 64),
		entries:  make([]Entry, 0),
	}
}

// Append serializes the event payload, links the previous block hash, and computes the SHA-256 digest.
func (c *Chain) Append(eventType string, event any) (Entry, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return Entry{}, err
	}

	entry := Entry{
		TimestampUnixMs: time.Now().UnixMilli(),
		EventType:       eventType,
		Event:           eventJSON,
		PrevHash:        c.prevHash,
	}

	// Compute hash of the pre-image
	data, err := json.Marshal(entry)
	if err != nil {
		return Entry{}, err
	}

	sum := sha256.Sum256(data)
	entry.Hash = hex.EncodeToString(sum[:])
	c.prevHash = entry.Hash
	c.entries = append(c.entries, entry)

	return entry, nil
}

// Entries returns a copy of all blocks in the ledger.
func (c *Chain) Entries() []Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res := make([]Entry, len(c.entries))
	copy(res, c.entries)
	return res
}

// Verify recomputes the cryptographic hash chain in forward order.
// Returns false immediately if any record was modified, reordered, or deleted post-hoc.
func Verify(entries []Entry) bool {
	prev := strings.Repeat("0", 64)

	for _, e := range entries {
		check := e
		check.Hash = "" // Clear hash field to reconstruct identical pre-image

		data, err := json.Marshal(check)
		if err != nil || check.PrevHash != prev {
			return false
		}

		sum := sha256.Sum256(data)
		computedHash := hex.EncodeToString(sum[:])
		if computedHash != e.Hash {
			return false
		}
		prev = e.Hash
	}

	return true
}
