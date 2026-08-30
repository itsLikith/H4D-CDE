package chain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

type Entry struct {
	TimestampUnixMs int64           `json:"ts"`
	Event           json.RawMessage `json:"event"`
	PrevHash        string          `json:"prev_hash"`
	Hash            string          `json:"hash,omitempty"`
}

type Chain struct{ prevHash string }

func New() *Chain { return &Chain{prevHash: strings.Repeat("0", 64)} }

func (c *Chain) Append(event any) (Entry, error) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return Entry{}, err
	}
	entry := Entry{TimestampUnixMs: time.Now().UnixMilli(), Event: eventJSON, PrevHash: c.prevHash}
	data, err := json.Marshal(entry)
	if err != nil {
		return Entry{}, err
	}
	sum := sha256.Sum256(data)
	entry.Hash = hex.EncodeToString(sum[:])
	c.prevHash = entry.Hash
	return entry, nil
}

// Verify recomputes every hash in sequence and returns false the moment any
// entry has been altered, reordered, or removed after the fact.
func Verify(entries []Entry) bool {
	prev := strings.Repeat("0", 64)
	for _, e := range entries {
		check := e
		check.Hash = ""
		data, err := json.Marshal(check)
		if err != nil || check.PrevHash != prev {
			return false
		}
		sum := sha256.Sum256(data)
		if hex.EncodeToString(sum[:]) != e.Hash {
			return false
		}
		prev = e.Hash
	}
	return true
}
