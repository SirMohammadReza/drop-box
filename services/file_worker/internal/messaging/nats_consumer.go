package messaging

import (
	"context"
	"encoding/json"
	"file_worker/internal/worker"
	"log"

	"github.com/nats-io/nats.go/jetstream"
)

type NatsConsumer struct {
	consumer  jetstream.Consumer
	processor UploadProcessor
}

type UploadProcessor interface {
	ProcessUpload(ctx context.Context, event worker.UploadedEvent) error
}

func NewNatsConsumer(consumer jetstream.Consumer, processor UploadProcessor) *NatsConsumer {
	return &NatsConsumer{consumer: consumer, processor: processor}
}

func (n *NatsConsumer) Start(ctx context.Context) (jetstream.ConsumeContext, error) {
	return n.consumer.Consume(func(msg jetstream.Msg) {
		var event worker.UploadedEvent
		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			log.Printf("malformed message, terminating (no retry): %v", err)
			msg.Term()
			return
		}

		if err := n.processor.ProcessUpload(ctx, event); err != nil {
			log.Printf("processing failed, will redeliver: %v", err)
			msg.Nak()
			return
		}

		msg.Ack()
	})
}
