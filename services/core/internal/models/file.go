package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type File struct {
	ID        bson.ObjectID  `bson:"_id,omitempty"`
	FolderID  *bson.ObjectID `bson:"folder_id"`
	UserId    uint           `bson:"user_id"`
	ObjectKey string         `bson:"object_key"`
	Name      string         `bson:"name"`
	Size      int64          `bson:"size"`
	MimeType  string         `bson:"mime_type"`
	CreatedAt time.Time      `bson:"created_at"`
	UpdatedAt time.Time      `bson:"updated_at"`
	DeletedAt time.Time      `bson:"deleted_at"`
}
