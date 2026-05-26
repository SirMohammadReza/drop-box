package uploader

import (
	"context"
	"io"
)

type Provider interface {
	UploadFile(c context.Context, metadata Metadata, data io.Reader) (*UploadResponse, error)
}

type Metadata struct {
	FileName string `json:"file_name"`
	UserUUID string `json:"user_uuid"`
	Size     int64  `json:"size"`
}

type UploadResponse struct {
	FileID  string `json:"file_id"`
	Success bool   `json:"success"`
}
