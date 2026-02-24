package folder

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Folder struct {
	ID        bson.ObjectID  `bson:"_id,omitempty"`
	UserID    uint           `bson:"user_id"`
	Name      string         `bson:"name"`
	ParentID  *bson.ObjectID `bson:"parent_id,omitempty"`
	Path      string         `bson:"path"`
	CreatedAt bson.Timestamp `bson:"created_at"`
	UpdatedAt bson.Timestamp `bson:"updated_at"`
}
