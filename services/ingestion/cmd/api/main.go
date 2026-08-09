package main

import (
	"context"
	"fmt"
	"ingestion/internal/config"
	"ingestion/internal/messaging"
	"ingestion/internal/platform/storage"
	"log"
	"net"
	"time"

	grpcHandlers "ingestion/internal/grpc"
	uploaderPb "ingestion/internal/proto/uploader"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/grpc"
)

var cfg config.Config

func main() {
	setEnv(&cfg)

	minioClient, err := minio.New(fmt.Sprintf("%s:%d", cfg.MinioURL, cfg.MinioPort), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioID, cfg.MinioSecret, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("failed to create minio client: %v", err)
	}
	fileStorage := storage.NewMinioStorage(minioClient, "drive-files")

	nc, err := nats.Connect(fmt.Sprintf("%s:%d", cfg.NatsURL, cfg.NatsPort))
	if err != nil {
		log.Fatalf("failed to connect to nats: %v", err)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("failed to init jetstream: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      "FILES",
		Subjects:  []string{"files.>"},
		Storage:   jetstream.FileStorage,
		Retention: jetstream.LimitsPolicy,
		MaxAge:    7 * 24 * time.Hour,
	})
	if err != nil {
		log.Fatalf("ensuring FILES stream: %v", err)
	}
	publisher := messaging.NewNatsPublisher(js, "files.uploaded")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen on %d: %v", cfg.GRPCPort, err)
	}

	s := grpc.NewServer()
	uplaoderHandler := grpcHandlers.NewUploaderHandler(fileStorage, publisher)
	uploaderPb.RegisterUploaderServiceServer(s, uplaoderHandler)

	log.Print("app started...")

	if err = s.Serve(lis); err != nil {
		log.Fatalf("could not run app: %v", err)
	}
}

func setEnv(cfg *config.Config) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("loading env: %v", err)
	}

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("parse env: %v", err)
	}
}
