package messaging

import (
	"context"
	"encoding/json"
	"file_worker/internal/worker"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type NatsPublisher struct {
	js      jetstream.JetStream
	subject string
}

func NewNatsPublisher(js jetstream.JetStream, subject string) *NatsPublisher {
	return &NatsPublisher{js: js, subject: subject}
}

func (p *NatsPublisher) PublishFileProcessed(ctx context.Context, event worker.ProcessedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshaling event: %w", err)
	}
	if _, err := p.js.Publish(ctx, p.subject, payload); err != nil {
		return fmt.Errorf("publishing to nats: %w", err)
	}
	return nil
}
