package main

import (
	"ingestion/internal/platform/storage"
	"log"
	"net"

	grpcHandlers "ingestion/internal/grpc"
	uploaderPb "ingestion/internal/proto/uploader"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen on 50052: %v", err)
	}

	s := grpc.NewServer()
	uplaoderHandler := grpcHandlers.NewUploaderHandler(fileStorage)
	uploaderPb.RegisterUploaderServiceServer(s, uplaoderHandler)

	log.Print("app started...")

	if err = s.Serve(lis); err != nil {
		log.Fatalf("could not run app: %v", err)
	}
}
