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

package ws

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// BridgeKafkaToClients consumes conflicts.detected events and broadcasts them to all connected operator dashboards.
func BridgeKafkaToClients(ctx context.Context, brokers []string, topic string, hub *Hub) error {
	if topic == "" {
		topic = "conflicts.detected"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "gateway-ws-group",
	})
	defer reader.Close()

	log.Printf("[*] Gateway WebSocket Kafka bridge consuming from '%s'", topic)

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

		hub.Broadcast(msg.Value)
	}
}