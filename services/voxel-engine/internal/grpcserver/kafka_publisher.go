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

package grpcserver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type EventPublisher interface {
	PublishAudit(ctx context.Context, event any)
	PublishConflictDetected(ctx context.Context, conflicts any)
	Close() error
}

type KafkaPublisher struct {
	auditWriter     *kafka.Writer
	conflictsWriter *kafka.Writer
}

func NewKafkaPublisher(brokers []string, auditTopic, conflictsTopic string) *KafkaPublisher {
	if auditTopic == "" {
		auditTopic = "audit.events"
	}
	if conflictsTopic == "" {
		conflictsTopic = "conflicts.detected"
	}

	return &KafkaPublisher{
		auditWriter: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    auditTopic,
			Balancer: &kafka.LeastBytes{},
			Async:    true,
		},
		conflictsWriter: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    conflictsTopic,
			Balancer: &kafka.LeastBytes{},
			Async:    true,
		},
	}
}

// PublishAudit sends an audit event asynchronously (fire-and-forget) so downstream logging latency
// never impacts the synchronous flight plan ingestion path.
func (p *KafkaPublisher) PublishAudit(ctx context.Context, event any) {
	if p == nil || p.auditWriter == nil {
		return
	}
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	go func() {
		if err := p.auditWriter.WriteMessages(context.Background(), kafka.Message{Value: data}); err != nil {
			log.Printf("[KafkaPublisher] audit publish warning: %v", err)
		}
	}()
}

// PublishConflictDetected broadcasts newly identified conflicts to Kafka for WebSocket fan-out to operator dashboards.
func (p *KafkaPublisher) PublishConflictDetected(ctx context.Context, conflicts any) {
	if p == nil || p.conflictsWriter == nil {
		return
	}
	data, err := json.Marshal(conflicts)
	if err != nil {
		return
	}
	go func() {
		if err := p.conflictsWriter.WriteMessages(context.Background(), kafka.Message{Value: data}); err != nil {
			log.Printf("[KafkaPublisher] conflict publish warning: %v", err)
		}
	}()
}

func (p *KafkaPublisher) Close() error {
	if p == nil {
		return nil
	}
	if p.auditWriter != nil {
		_ = p.auditWriter.Close()
	}
	if p.conflictsWriter != nil {
		_ = p.conflictsWriter.Close()
	}
	return nil
}

// NoOpPublisher for standalone execution and unit testing without active Kafka broker.
type NoOpPublisher struct{}

func (n *NoOpPublisher) PublishAudit(_ context.Context, _ any)            {}
func (n *NoOpPublisher) PublishConflictDetected(_ context.Context, _ any) {}
func (n *NoOpPublisher) Close() error                                     { return nil }
