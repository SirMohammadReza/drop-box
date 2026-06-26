package grpc

import (
	"bufio"
	"fmt"
	"ingestion/internal/proto/uploader"
	"io"
	"os"
	"path/filepath"
)

type UploaderHandler struct {
	uploader.UnimplementedUploaderServiceServer
	storageDir string
}

func NewUploaderHandler(sd string) *UploaderHandler {
	return &UploaderHandler{
		storageDir: sd,
	}
}

func (uh *UploaderHandler) UploadFile(stream uploader.UploaderService_UploadFileServer) error {
	pr, pw := io.Pipe()

	firstReq, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("could not get request: %v", err)
	}

	metadata := firstReq.GetMetadata()
	if metadata == nil {
		return fmt.Errorf("first message must be metadata: %v", err)
	}

	fileExt := filepath.Ext(metadata.FileName)
	fileID := metadata.FileId
	dst := filepath.Join(uh.storageDir, fileID+fileExt)

	go func() {
		defer pw.Close()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			if chunk := req.GetChunk(); chunk != nil {
				if _, err := pw.Write(chunk); err != nil {
					return
				}
			}
		}
	}()

	file, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	bufWriter := bufio.NewWriter(file)

	_, err = io.Copy(bufWriter, pr)
	if err != nil {
		return fmt.Errorf("could not save file: %v", err)
	}

	if err := bufWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %v", err)
	}

	return stream.SendAndClose(&uploader.UploadResponse{
		FileId:  fileID,
		Success: true,
	})
}
