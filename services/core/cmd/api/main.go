package main

import (
	"log"

	"core/internal/auth"
	authHttp "core/internal/auth/delivery/http"
	tokenPb "core/internal/proto/token"
	uploaderPb "core/internal/proto/uploader"
	userPb "core/internal/proto/user"
	"core/internal/uploader"
	uploaderHttp "core/internal/uploader/delivery/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to grpc: %v", err)
	}
	defer grpcConn.Close()

	userClient := userPb.NewUserServiceClient(grpcConn)
	tokenClient := tokenPb.NewTokenServiceClient(grpcConn)
	uploaderClient := uploaderPb.NewUploaderServiceClient(grpcConn)

	authProvider := auth.NewGRPCAdapter(userClient, tokenClient)
	authHandler := authHttp.NewHandler(authProvider)

	uploaderProvider := uploader.NewGRPCAdapter(uploaderClient)
	uploaderHandler := uploaderHttp.NewHandler(uploaderProvider)

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	api := e.Group("/api")
	v1 := api.Group("/v1")
	apiGroup := v1.Group("/auth")
	authHandler.RegisterRoutes(apiGroup)

	apiGroup = v1.Group("/files")
	uploaderHandler.RegisterRoutes(apiGroup)

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("could not start app: %v", err)
	}
}
