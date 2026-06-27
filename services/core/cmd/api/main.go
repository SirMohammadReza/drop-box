package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"core/internal/auth"
	authHttp "core/internal/auth/delivery/http"
	"core/internal/file"
	mongoDB "core/internal/platform/mongo"
	tokenPb "core/internal/proto/token"
	uploaderPb "core/internal/proto/uploader"
	userPb "core/internal/proto/user"
	"core/internal/uploader"
	uploaderHttp "core/internal/uploader/delivery/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mongoClient := mongoDB.ConnetMongo()

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	v1 := e.Group("/api/v1")

	authConn := authInit(v1)
	defer authConn.Close()

	ingestionConn := ingestionInit(v1, mongoClient)
	defer ingestionConn.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         ":8080",
		GracefulTimeout: 10 * time.Second,
	}

	if err := sc.Start(ctx, e); err != nil {
		log.Fatal("server error:", err)
	}
}

func ingestionInit(echoGroup *echo.Group, mc *mongo.Client) *grpc.ClientConn {
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to ingestion grpc: %v", err)
	}

	uploaderClient := uploaderPb.NewUploaderServiceClient(conn)

	fileRepo := file.NewMongoRepository(mc)
	fileService := file.NewFileService(fileRepo)

	uploaderProvider := uploader.NewGRPCAdapter(uploaderClient)
	uploaderHandler := uploaderHttp.NewHandler(uploaderProvider, fileService)

	uploaderHandler.RegisterRoutes(echoGroup.Group("/files"))

	return conn
}

func authInit(echoGroup *echo.Group) *grpc.ClientConn {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to authenticationgrpc: %v", err)
	}

	userClient := userPb.NewUserServiceClient(conn)
	tokenClient := tokenPb.NewTokenServiceClient(conn)

	authProvider := auth.NewGRPCAdapter(userClient, tokenClient)
	authHandler := authHttp.NewHandler(authProvider)

	authHandler.RegisterRoutes(echoGroup.Group("/auth"))

	return conn
}
