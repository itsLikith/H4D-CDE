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

package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/itsLikith/h4d-cde/services/audit-svc/internal/chain"
	"github.com/segmentio/kafka-go"
)

// Run continuously consumes audit messages from the Kafka audit.events stream,
// hash-chains them in exact arrival order, and passes them to the storage layer.
func Run(
	ctx context.Context,
	brokers []string,
	topic string,
	groupID string,
	c *chain.Chain,
	persist func(chain.Entry) error,
) error {
	if topic == "" {
		topic = "audit.events"
	}
	if groupID == "" {
		groupID = "audit-svc"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer reader.Close()

	log.Printf("[*] Audit consumer listening on topic '%s' (group: '%s')", topic, groupID)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var rawMap map[string]any
		if err := json.Unmarshal(msg.Value, &rawMap); err != nil {
			log.Printf("[!] Skipping malformed audit payload: %v", err)
			continue
		}

		eventType, _ := rawMap["type"].(string)
		if eventType == "" {
			eventType = "audit_event"
		}

		entry, err := c.Append(eventType, rawMap)
		if err != nil {
			log.Printf("[!] Hash chain append failure: %v", err)
			continue
		}

		if persist != nil {
			if err := persist(entry); err != nil {
				log.Printf("[!] Storage persistence warning: %v", err)
			}
		}
	}
}
