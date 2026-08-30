package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"hive/audit-svc/internal/chain"
)

// Run consumes every event published to audit.events by any other service
// and hash-chains them in arrival order.
func Run(ctx context.Context, brokers []string, c *chain.Chain, persist func(chain.Entry) error) error {
	reader := kafka.NewReader(kafka.ReaderConfig{Brokers: brokers, Topic: "audit.events", GroupID: "audit-svc"})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		var event map[string]any
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("skipping malformed audit event: %v", err)
			continue
		}
		entry, err := c.Append(event)
		if err != nil {
			return err
		}
		if err := persist(entry); err != nil {
			return err
		}
	}
}
