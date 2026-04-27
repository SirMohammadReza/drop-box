package logger

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoHook struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoHook(client *mongo.Client, dbName, collectionName string) *MongoHook {
	return &MongoHook{
		client:     client,
		collection: client.Database(dbName).Collection(collectionName),
	}
}

func (mh *MongoHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (mh *MongoHook) Fire(entry *logrus.Entry) error {
	data := make(logrus.Fields)
	for k, v := range entry.Data {
		data[k] = v
	}

	data["level"] = entry.Level.String()
	data["message"] = entry.Message
	data["time"] = entry.Time

	go func(logData logrus.Fields) {
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, _ = mh.collection.InsertOne(c, logData)
	}(data)

	return nil
}
