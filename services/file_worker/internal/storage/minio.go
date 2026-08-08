package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioClient(client *minio.Client, bucket string) *MinioStorage {
	return &MinioStorage{client: client, bucket: bucket}
}

func (m *MinioStorage) Get(ctx context.Context, name string) (io.ReadCloser, error) {
	return m.client.GetObject(ctx, m.bucket, name, minio.GetObjectOptions{})
}

func (m *MinioStorage) Create(ctx context.Context, name string) (io.WriteCloser, error) {
	pr, pw := io.Pipe()
	done := make(chan error, 1)
	go func() {
		_, err := m.client.PutObject(ctx, m.bucket, name, pr, -1, minio.PutObjectOptions{})
		pr.CloseWithError(err)
		done <- err
	}()

	return &pipeWriteCloser{
		done:       done,
		PipeWriter: pw,
	}, nil
}

func (m *MinioStorage) Remove(ctx context.Context, name string) error {
	return m.client.RemoveObject(ctx, m.bucket, name, minio.RemoveObjectOptions{})
}

type pipeWriteCloser struct {
	*io.PipeWriter
	done chan error
}

func (w *pipeWriteCloser) Close() error {
	if err := w.PipeWriter.Close(); err != nil {
		return err
	}
	return <-w.done
}
