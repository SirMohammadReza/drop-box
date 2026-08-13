package status

import (
	"context"
	"core/internal/file"
)

type FileUploadedEvent struct {
	FileID    string `json:"file_id"`
	ObjectKey string `json:"object_key"`
	FileName  string `json:"file_name"`
	UserUuid  string `json:"user_uuid"`
	Size      uint64 `json:"size"`
}

type FileRepository interface {
	UpdateStatus(ctx context.Context, fileID string, status file.FileStatus) error
}

type ZippingStatusConsumer struct {
	repo FileRepository
}

func NewZippingStatusConsumer(repo FileRepository) *ZippingStatusConsumer {
	return &ZippingStatusConsumer{repo: repo}
}

func (s *ZippingStatusConsumer) HandleFileUpload(ctx context.Context, event FileUploadedEvent) error {
	return s.repo.UpdateStatus(ctx, event.FileID, file.FileStatusZipping)
}
