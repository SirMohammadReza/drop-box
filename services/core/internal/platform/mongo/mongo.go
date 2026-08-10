package mongo

import (
	"context"
	"core/internal/config"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	mongoClient *mongo.Client
	once        sync.Once
)

func ConnetMongo(cfg *config.Config) *mongo.Client {
	once.Do(func() {
		uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.MongoUsername, cfg.MongoPassword, cfg.MongoURL, cfg.MongoPort)

		clientOption := options.Client().ApplyURI(uri)

		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clien, err := mongo.Connect(clientOption)
		if err != nil {
			log.Fatalf("Error while connecting mongo: %v", err.Error())
		}

		err = clien.Ping(c, nil)
		if err != nil {
			log.Fatalf("Could not connect to mongo: %v", err)
		}

		mongoClient = clien
	})

	return mongoClient
}
