package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("could not load env: %v", err)
	}

	return os.Getenv(key)
}
