package grpc

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"ingestion/internal/proto/uploader"
	"io"
	"path"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxFileNameLen = 255

type FileStorage interface {
	Create(ctx context.Context, name string) (io.WriteCloser, error)
	Remove(ctx context.Context, name string) error
}
type UploaderHandler struct {
	uploader.UnimplementedUploaderServiceServer
	storage   FileStorage
	publisher EventPublisher
}

type FileUploadEvent struct {
	FileID    string `json:"file_id"`
	ObjectKet string `json:"object_key"`
	FileName  string `json:"file_name"`
	UserUuid  string `json:"user_uuid"`
	Size      uint64 `json:"size"`
}

type EventPublisher interface {
	PublishFileUploaded(ctx context.Context, event FileUploadEvent) error
}

func NewUploaderHandler(storage FileStorage, publisher EventPublisher) *UploaderHandler {
	return &UploaderHandler{storage: storage, publisher: publisher}
}

func (h *UploaderHandler) UploadFile(stream uploader.UploaderService_UploadFileServer) error {
	ctx := stream.Context()

	firstReq, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Internal, "reading first message: %v", err)
	}

	metadata := firstReq.GetMetadata()
	if metadata == nil {
		return status.Error(codes.InvalidArgument, "first message must contains metadata")
	}

	name, err := objectName(metadata)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid metadata: %v", err)
	}

	w, err := h.storage.Create(ctx, name)
	if err != nil {
		return status.Errorf(codes.Internal, "creating destination object: %v", err)
	}

	if err := writeChunks(stream, w); err != nil {
		w.Close()
		h.storage.Remove(ctx, name)
		return err
	}

	if err := w.Close(); err != nil {
		h.storage.Remove(ctx, name)
		return status.Errorf(codes.Internal, "closing object: %v", err)
	}

	event := FileUploadEvent{
		FileID:    metadata.FileId,
		ObjectKet: name,
		FileName:  metadata.FileName,
		UserUuid:  metadata.UserUuid,
	}
	if err := h.publisher.PublishFileUploaded(ctx, event); err != nil {
		h.storage.Remove(ctx, name)
		return status.Errorf(codes.Internal, "publishing upload event: %v", err)
	}

	return stream.SendAndClose(&uploader.UploadResponse{
		FileId:  metadata.FileId,
		Success: true,
	})
}

func objectName(metadata *uploader.Metadata) (string, error) {
	if metadata.FileId == "" {
		return "", errors.New("file_id is required")
	}
	if metadata.FileName == "" {
		return "", errors.New("file_name is required")
	}
	if len(metadata.FileName) > maxFileNameLen {
		return "", fmt.Errorf("file_name exceeds %d characters", maxFileNameLen)
	}
	if metadata.FileId != path.Base(metadata.FileId) || strings.Contains(metadata.FileId, "..") {
		return "", errors.New("file_id must not contain path separators")
	}

	ext := path.Ext(metadata.FileName)
	return "raw/" + metadata.FileId + ext, nil
}

func writeChunks(stream uploader.UploaderService_UploadFileServer, w io.Writer) error {
	bufWrite := bufio.NewWriter(w)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "reciving chunk: %v", err)
		}
		if chunk := req.GetChunk(); len(chunk) > 0 {
			if _, err := bufWrite.Write(chunk); err != nil {
				return status.Errorf(codes.Internal, "writing chunk: %v", err)
			}
		}
	}

	if err := bufWrite.Flush(); err != nil {
		return status.Errorf(codes.Internal, "flushing file: %v", err)
	}

	return nil
}
