package main

import (
	"log"

	"core/internal/auth"
	"core/internal/auth/delivery/http"
	tokenPb "core/internal/proto/token"
	userPb "core/internal/proto/user"

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

	userClient := userPb.NewUserServiceClient(grpcConn)
	tokenClient := tokenPb.NewTokenServiceClient(grpcConn)

	authProvider := auth.NewGRPCAdapter(userClient, tokenClient)
	authHandler := http.NewHandler(authProvider)

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	apiGroup := e.Group("api/v1/auth")
	authHandler.RegisterRoutes(apiGroup)

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("could not start app: %v", err)
	}
}
