package messaging

import (
	"context"
	"core/internal/status"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go/jetstream"
)

type UploadHandler interface {
	HandleFileUpload(ctx context.Context, event status.FileUploadedEvent) error
}

type NatsConsumer struct {
	consumer jetstream.Consumer
	handler  UploadHandler
}

func NewNatsConsumer(consumer jetstream.Consumer, handler UploadHandler) *NatsConsumer {
	return &NatsConsumer{consumer: consumer, handler: handler}
}

func (n *NatsConsumer) Start(ctx context.Context) (jetstream.ConsumeContext, error) {
	return n.consumer.Consume(func(msg jetstream.Msg) {
		var event status.FileUploadedEvent
		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			log.Printf("malformed message, terminating (no retry): %v", err)
			msg.Term()
			return
		}

		if err := n.handler.HandleFileUpload(ctx, event); err != nil {
			log.Printf("updating status failed, will redeliver: %v", err)
			msg.Nak()
			return
		}

		msg.Ack()
	})
}
