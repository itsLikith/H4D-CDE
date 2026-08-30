// services/gateway/internal/ws/kafka_bridge.go
package ws

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// BridgeKafkaToClients consumes conflicts.detected (published by
// voxel-engine, Part 4.4) and fans each event out to every connected
// dashboard -- this is the async half of the Part 4.3 sequence diagram.
func BridgeKafkaToClients(ctx context.Context, brokers []string, hub *Hub) error {
	reader := kafka.NewReader(kafka.ReaderConfig{Brokers: brokers, Topic: "conflicts.detected", GroupID: "gateway-ws"})
	defer reader.Close()
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		hub.Broadcast(msg.Value)
	}
}