package main

import (
	"ingestion/internal/messaging"
	"ingestion/internal/platform/storage"
	"log"
	"net"

	grpcHandlers "ingestion/internal/grpc"
	uploaderPb "ingestion/internal/proto/uploader"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/grpc"
)

func main() {
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("admin", "password123", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("failed to create minio client: %v", err)
	}
	fileStorage := storage.NewMinioStorage(minioClient, "drive-files")

	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("failed to connect to nats: %v", err)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("failed to init jetstream: %v", err)
	}
	publisher := messaging.NewNatsPublisher(js, "files.uploaded")

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen on 50052: %v", err)
	}

	s := grpc.NewServer()
	uplaoderHandler := grpcHandlers.NewUploaderHandler(fileStorage, publisher)
	uploaderPb.RegisterUploaderServiceServer(s, uplaoderHandler)

	log.Print("app started...")

	if err = s.Serve(lis); err != nil {
		log.Fatalf("could not run app: %v", err)
	}
}
