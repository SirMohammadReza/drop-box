package main

import (
	grpcHandlers "authentication/internal/grpc"
	"authentication/internal/platform/db/postgres"
	"authentication/internal/platform/encryption"
	"authentication/internal/token"
	tokenProto "authentication/internal/token/proto/token"
	"authentication/internal/user"
	userProto "authentication/internal/user/proto/user"
	"errors"
	"log"
	"net"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
)

func main() {
	postgresDB := postgres.GetDB()
	// runMigrations()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pm := encryption.NewEncryptionService()

	tr := token.NewRepository(postgresDB)
	ts := token.NewTokenService(tr)

	ur := user.NewRepository(postgresDB)
	us := user.NewUserService(ur, pm, ts)

	userHandler := grpcHandlers.NewUserHandler(us)
	userProto.RegisterUserServiceServer(s, userHandler)

	tokenHandler := grpcHandlers.NewTokenHandler(ts)
	tokenProto.RegisterTokenServiceServer(s, tokenHandler)

	log.Println("App started...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve grpc server: %v", err)
	}

}

func runMigrations() {
	dsn := "postgres://admin:4321@postgres:5432/auth?sslmode=disable"
	migrations, err := migrate.New("file://internal/migrations", dsn)
	if err != nil {
		log.Fatalf("Could not run migrations instance: %v", err)
	}

	if err = migrations.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Make no change to databases")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		log.Println("Migrations run successfully")
	}
}
