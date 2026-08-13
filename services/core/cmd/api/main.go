package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"core/internal/auth"
	authHttp "core/internal/auth/delivery/http"
	"core/internal/config"
	"core/internal/file"
	"core/internal/platform/messaging"
	mongoDB "core/internal/platform/mongo"
	tokenPb "core/internal/proto/token"
	uploaderPb "core/internal/proto/uploader"
	userPb "core/internal/proto/user"
	"core/internal/status"
	"core/internal/uploader"
	uploaderHttp "core/internal/uploader/delivery/http"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cfg config.Config

func main() {
	setEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mongoClient := mongoDB.ConnetMongo(&cfg)

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	v1 := e.Group("/api/v1")

	authConn := authInit(v1)
	defer authConn.Close()

	ingestionConn := ingestionInit(v1, mongoClient)
	defer ingestionConn.Close()

	fileRepo := file.NewMongoRepository(mongoClient)

	natsConn, consumectx := natsInit(ctx, fileRepo)
	defer natsConn.Close()
	defer consumectx.Stop()

	sc := echo.StartConfig{
		Address:         ":8080",
		GracefulTimeout: 10 * time.Second,
	}

	if err := sc.Start(ctx, e); err != nil {
		log.Fatal("server error:", err)
	}
}

func ingestionInit(echoGroup *echo.Group, mc *mongo.Client) *grpc.ClientConn {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.GRPCURL, cfg.GRPCIngestionPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.GRPCURL, cfg.GRPCAuthPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func setEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("loading env: %v", err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("parsing env: %v", err)
	}
}

func natsInit(ctx context.Context, fileRepo *file.MongoRepository) (*nats.Conn, jetstream.ConsumeContext) {
	nc, err := nats.Connect(fmt.Sprintf("%s:%d", cfg.NatsURL, cfg.NatsPort))
	if err != nil {
		log.Fatalf("could not connect to nats: %v", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("could not init jetstream: %v", err)
	}

	stream, err := js.Stream(ctx, "FILES")
	if err != nil {
		log.Fatalf("stream lookup (is ingestion running?): %v", err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       "core-status-updater",
		FilterSubject: "files.uploaded",
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatalf("ensuring nats consumer: %v", err)
	}

	statusConsumer := status.NewZippingStatusConsumer(fileRepo)
	natsConsumer := messaging.NewNatsConsumer(consumer, statusConsumer)

	consumeCtx, err := natsConsumer.Start(ctx)
	if err != nil {
		log.Fatalf("starting nats consumer: %v", err)
	}

	return nc, consumeCtx
}
