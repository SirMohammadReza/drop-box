package grpc

import (
	"bufio"
	"errors"
	"fmt"
	"ingestion/internal/proto/uploader"
	"io"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxFileNameLen = 255

type UploaderHandler struct {
	uploader.UnimplementedUploaderServiceServer
	storageDir string
}

func NewUploaderHandler(storageDir string) *UploaderHandler {
	return &UploaderHandler{storageDir: storageDir}
}

func (h *UploaderHandler) UploadFile(stream uploader.UploaderService_UploadFileServer) error {
	firstReq, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Internal, "reading first message: %v", err)
	}

	metadata := firstReq.GetMetadata()
	if metadata == nil {
		return status.Error(codes.InvalidArgument, "first message must contains metadata")
	}

	dst, err := h.destinationPath(metadata)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid metadata: %v", err)
	}

	file, err := os.Create(dst)
	if err != nil {
		return status.Errorf(codes.Internal, "creating destination file: %v", err)
	}

	if err := writeChunks(stream, file); err != nil {
		file.Close()
		os.Remove(dst)
		return err
	}

	if err := file.Close(); err != nil {
		os.Remove(dst)
		return status.Errorf(codes.Internal, "closing file: %v", err)
	}

	return stream.SendAndClose(&uploader.UploadResponse{
		FileId:  metadata.FileId,
		Success: true,
	})
}

func (h *UploaderHandler) destinationPath(metadata *uploader.Metadata) (string, error) {
	if metadata.FileId == "" {
		return "", errors.New("file_id is required")
	}
	if metadata.FileName == "" {
		return "", errors.New("file_name is required")
	}
	if len(metadata.FileName) > maxFileNameLen {
		return "", fmt.Errorf("file_name exceeds %d characters", maxFileNameLen)
	}
	if filepath.Base(metadata.FileId) != metadata.FileId {
		return "", errors.New("file_id must not contain path separators")
	}

	ext := filepath.Ext(metadata.FileName)
	dst := filepath.Join(h.storageDir, metadata.FileId+ext)
	if !isSubPath(h.storageDir, dst) {
		return "", errors.New("resolved path escapes storage directory")
	}

	return dst, nil
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

func isSubPath(base, target string) bool {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
