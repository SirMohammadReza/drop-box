package main

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// postgresDB := postgres.GetDB()
	dsn := "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"
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
