package file

import "go.mongodb.org/mongo-driver/v2/bson"

type NewFileInputs struct {
	FolderID  *bson.ObjectID
	UserID    uint
	ObjectKey string
	Name      string
	Size      int64
	MimeType  string
}
