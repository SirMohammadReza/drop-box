package mongo

import (
	"authentication/internal/config"
	"context"
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

func ConnetMongo() *mongo.Client {
	once.Do(func() {
		uri := getUri()

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

func getUri() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", config.Config("MONGO_USERNAME"), config.Config("MONGO_PASSWORD"), config.Config("MONGO_HOST"), config.Config("MONGO_PORT"))
}
