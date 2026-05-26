package main

import (
	grpcHandlers "ingestion/internal/grpc"
	uploaderPb "ingestion/internal/proto/uploader"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen on 50052: %v", err)
	}

	s := grpc.NewServer()
	uplaoderHandler := grpcHandlers.NewUploaderHandler("./files/")
	uploaderPb.RegisterUploaderServiceServer(s, uplaoderHandler)

	log.Print("app started...")

	if err = s.Serve(lis); err != nil {
		log.Fatalf("could not run app: %v", err)
	}
}
