package main

import (
	"context"
	"file_worker/internal/config"
	"file_worker/internal/messaging"
	"file_worker/internal/storage"
	"file_worker/internal/worker"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var cfg config.Config

func main() {
	setEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	minioClient, err := minio.New(fmt.Sprintf("%s:%d", cfg.MinioURL, cfg.MinioPort), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioID, cfg.MinioSecret, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("minio client: %v", err)
	}
	fileStorage := storage.NewMinioClient(minioClient, "drive-files")

	nc, err := nats.Connect(fmt.Sprintf("%s:%d", cfg.NatsURL, cfg.NatsPort))
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

	stream, err := js.Stream(ctx, "FILES")
	if err != nil {
		log.Fatalf("stream lookup: %v", err)
	}
	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       "file-worker",
		FilterSubject: "files.uploaded",
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatalf("ensuring consumer: %v", err)
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

func setEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("loading env: %v", err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("parse env: %v", err)
	}
}
