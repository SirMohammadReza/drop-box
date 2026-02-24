package folder

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(mc *mongo.Client) *MongoRepository {
	return &MongoRepository{
		collection: mc.Database("core").Collection("folders"),
	}
}

func (mr *MongoRepository) Create(c context.Context, f *Folder) error {
	if f.ParentID != nil {
		var parent Folder
		err := mr.collection.FindOne(c, bson.M{"_id": f.ParentID}).Decode(&parent)
		if err != nil {
			return fmt.Errorf("failed to find parent: %w", err)
		}
		f.Path = fmt.Sprintf("%s%s/", parent.Path, parent.ID.Hex())
	} else {
		f.Path = "/"
	}

	f.ID = bson.NewObjectID()
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()

	_, err := mr.collection.InsertOne(c, f)
	return err
}

func (mr *MongoRepository) GetFolderByName(c context.Context, folderName string) (*Folder, error) {
	var folder Folder
	filter := bson.M{"name": folderName}

	if err := mr.collection.FindOne(c, filter).Decode(&folder); err != nil {
		return nil, err
	}

	return &folder, nil
}
