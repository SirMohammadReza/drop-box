package main

import (
	"context"
	"file_worker/internal/messaging"
	"file_worker/internal/storage"
	"file_worker/internal/worker"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("admin", "password123", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("minio client: %v", err)
	}
	fileStorage := storage.NewMinioClient(minioClient, "drive-files")

	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer nc.Close()
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("jetstream: %v", err)
	}

	publisher := messaging.NewNatsPublisher(js, "files.processed")
	processor := worker.NewFileprocessor(fileStorage, publisher)

	stream, err := js.Stream(ctx, "UPLOADS")
	if err != nil {
		log.Fatalf("stream lookup: %v", err)
	}

	consumer, err := stream.Consumer(ctx, "file-worker")
	if err != nil {
		log.Fatalf("consumer lookup: %v", err)
	}

	natsConsumer := messaging.NewNatsConsumer(consumer, processor)
	consumeCtx, err := natsConsumer.Start(ctx)
	if err != nil {
		log.Fatalf("starting consumer: %v", err)
	}
	defer consumeCtx.Stop()

	log.Print("file-worker started...")
	<-ctx.Done()
	log.Print("shutting down...")
}
