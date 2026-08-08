package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"ingestion/internal/grpc"

	"github.com/nats-io/nats.go/jetstream"
)

type NatsPublisher struct {
	js      jetstream.JetStream
	subject string
}

func NewNatsPublisher(js jetstream.JetStream, subject string) *NatsPublisher {
	return &NatsPublisher{js: js, subject: subject}
}

func (p *NatsPublisher) PublishFileUploaded(ctx context.Context, event grpc.FileUploadEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshaling event: %w", err)
	}

	if _, err := p.js.Publish(ctx, p.subject, payload); err != nil {
		return fmt.Errorf("publishing to nats: %w", err)
	}

	return nil
}
