package file

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(mc *mongo.Client) *MongoRepository {
	return &MongoRepository{
		collection: mc.Database("core").Collection("files"),
	}
}

func (mr *MongoRepository) Create(c context.Context, f *File) error {
	f.ID = bson.NewObjectID()
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()

	_, err := mr.collection.InsertOne(c, f)
	return err
}

func (mr *MongoRepository) GetFilesByFolder(c context.Context, folderID *bson.ObjectID) ([]File, error) {
	filter := bson.M{"folder_id": folderID}

	cursor, err := mr.collection.Find(c, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var files []File
	if err := cursor.All(c, &files); err != nil {
		return nil, err
	}

	return files, nil
}

func (mr *MongoRepository) Move(c context.Context, fileID *bson.ObjectID, newFolderID *bson.ObjectID) error {
	filter := bson.M{"_id": fileID}
	update := bson.M{"$set": bson.M{"folder_id": newFolderID, "updated_at": time.Now()}}

	_, err := mr.collection.UpdateOne(c, filter, update)
	return err
}

func (mr *MongoRepository) Copy(c context.Context, f File, destinationFolderID *bson.ObjectID) error {
	f.ID = bson.NewObjectID()
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()

	_, err := mr.collection.InsertOne(c, f)
	return err
}
