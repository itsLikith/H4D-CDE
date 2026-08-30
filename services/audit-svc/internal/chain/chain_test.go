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

package chain_test

import (
	"encoding/json"
	"testing"

	"github.com/itsLikith/h4d-cde/services/audit-svc/internal/chain"
	"github.com/stretchr/testify/assert"
)

func TestAuditChainIntegrityAndTamperDetection(t *testing.T) {
	c := chain.New()

	e1, err := c.Append("voxel_write", map[string]any{"entity": "UAV-1", "key": "8828308281fffff:400:100"})
	assert.NoError(t, err)
	assert.NotEmpty(t, e1.Hash)

	e2, err := c.Append("conflicts_detected", map[string]any{"conflicts": 2})
	assert.NoError(t, err)
	assert.Equal(t, e1.Hash, e2.PrevHash)

	e3, err := c.Append("advisory_applied", map[string]any{"strategy": "delay", "delay_s": 60})
	assert.NoError(t, err)
	assert.Equal(t, e2.Hash, e3.PrevHash)

	entries := c.Entries()
	assert.Len(t, entries, 3)
	assert.True(t, chain.Verify(entries), "Cryptographic chain must verify as valid")

	// Tampering test: modify an interior block payload
	tamperedEntries := make([]chain.Entry, len(entries))
	copy(tamperedEntries, entries)
	tamperedEntries[1].Event = json.RawMessage(`{"conflicts": 0, "tampered": true}`)

	assert.False(t, chain.Verify(tamperedEntries), "Chain must fail verification if an event payload was altered")

	// Tampering test: delete an interior block
	deletedEntries := []chain.Entry{entries[0], entries[2]}
	assert.False(t, chain.Verify(deletedEntries), "Chain must fail verification if an event was deleted")
}
