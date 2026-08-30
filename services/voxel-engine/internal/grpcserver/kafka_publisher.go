// services/voxel-engine/internal/grpcserver/kafka_publisher.go
package grpcserver

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	auditWriter     *kafka.Writer
	conflictsWriter *kafka.Writer
}

func (p *KafkaPublisher) PublishAudit(ctx context.Context, event any) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	// Fire-and-forget: run in a goroutine so a slow/unavailable Kafka broker
	// never adds latency to IngestFlightPlan's response (Part 4.5's NFR).
	go func() { _ = p.auditWriter.WriteMessages(ctx, kafka.Message{Value: data}) }()
}