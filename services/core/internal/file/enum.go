package file

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/x/bsonx/bsoncore"
)

type FileStatus int

const (
	FileStatusUnknown FileStatus = iota
	FileStatusUploading
	FileStatusZipping
	FileStatusDone
	FileStatusPending
)

var fileStatusStrings = map[FileStatus]string{
	FileStatusUnknown:   "UNKNOWN",
	FileStatusUploading: "UPLOADING",
	FileStatusZipping:   "ZIPPING",
	FileStatusDone:      "DONE",
	FileStatusPending:   "PENDING",
}

var fileStatusFromString = func() map[string]FileStatus {
	m := make(map[string]FileStatus, len(fileStatusStrings))
	for k, v := range fileStatusStrings {
		m[v] = k
	}
	return m
}()

func (f FileStatus) String() string {
	if s, ok := fileStatusStrings[f]; ok {
		return s
	}
	return "UNKNOWN"
}

func (f FileStatus) IsValid() bool {
	_, ok := fileStatusStrings[f]
	return ok
}

func (f FileStatus) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%q", f.String()), nil
}

func (f *FileStatus) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("FileStatus must be a string")
	}
	s := string(data[1 : len(data)-1])
	return f.fromString(s)
}

func (f FileStatus) MarshalBSONValue() (bson.Type, []byte, error) {
	return bson.TypeString, bsoncore.AppendString(nil, f.String()), nil
}

func (f *FileStatus) UnmarshalBSONValue(t bson.Type, data []byte) error {
	if t != bson.TypeString {
		return fmt.Errorf("cannot unmarshal BSON type %v into FileStatus", t)
	}
	s, _, ok := bsoncore.ReadString(data)
	if !ok {
		return fmt.Errorf("failed to read BSON string data")
	}
	return f.fromString(s)
}

func (f *FileStatus) fromString(s string) error {
	if status, ok := fileStatusFromString[s]; ok {
		*f = status
		return nil
	}
	*f = FileStatusUnknown
	return fmt.Errorf("unknown FileStatus %q", s)
}
