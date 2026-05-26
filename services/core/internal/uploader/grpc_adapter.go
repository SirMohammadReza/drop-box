package uploader

import (
	"context"
	uploaderPb "core/internal/proto/uploader"
	"io"
)

type grpcAdapter struct {
	uploaderClient uploaderPb.UploaderServiceClient
}

func NewGRPCAdapter(uc uploaderPb.UploaderServiceClient) Provider {
	return &grpcAdapter{
		uploaderClient: uc,
	}
}

func (a *grpcAdapter) UploadFile(c context.Context, meta Metadata, data io.Reader) (*UploadResponse, error) {
	stream, err := a.uploaderClient.UploadFile(c)
	if err != nil {
		return nil, err
	}

	err = stream.Send(&uploaderPb.UploadRequest{
		Data: &uploaderPb.UploadRequest_Metadata{
			Metadata: &uploaderPb.Metadata{
				FileName: meta.FileName,
				UserUuid: meta.UserUUID,
				Size:     meta.Size,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 32*1024)
	for {
		n, err := data.Read(buffer)
		if n > 0 {
			err := stream.Send(&uploaderPb.UploadRequest{
				Data: &uploaderPb.UploadRequest_Chunk{
					Chunk: buffer[:n],
				},
			})
			if err != nil {
				return nil, err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return &UploadResponse{
		FileID:  res.FileId,
		Success: res.Success,
	}, nil
}
