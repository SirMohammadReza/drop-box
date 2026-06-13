package file

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type File struct {
	ID        bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	Status    FileStatus     `bson:"status" json:"status"`
	FolderID  *bson.ObjectID `bson:"folder_id" json:"folder_id"`
	UserID    uint           `bson:"user_id" json:"user_id"`
	ObjectKey string         `bson:"object_key" json:"object_key"`
	Name      string         `bson:"name" json:"name"`
	Size      int64          `bson:"size" json:"size"`
	MimeType  string         `bson:"mime_type" json:"mime_type"`
	CreatedAt time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time      `bson:"updated_at" json:"updated_at"`
	DeletedAt time.Time      `bson:"deleted_at" json:"deleted_at"`
}
