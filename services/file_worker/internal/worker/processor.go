package worker

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

type UploadedEvent struct {
	FileID    string `json:"file_id"`
	ObjectKey string `json:"object_key"`
	FileName  string `json:"file_name"`
	UserUuid  string `json:"user_uuid"`
	Size      int64  `json:"size"`
}

type ProcessedEvent struct {
	FileId       string
	ProcessedKey string
	Checksum     string
}

type FileStorage interface {
	Get(ctx context.Context, name string) (io.ReadCloser, error)
	Create(ctx context.Context, name string) (io.WriteCloser, error)
	Remove(ctx context.Context, name string) error
}

type EventPublisher interface {
	PublishFileProcessed(ctx context.Context, event ProcessedEvent) error
}

type FileProcessor struct {
	storage   FileStorage
	publisher EventPublisher
}

func NewFileprocessor(storage FileStorage, publisher EventPublisher) *FileProcessor {
	return &FileProcessor{storage: storage, publisher: publisher}
}

func (fp *FileProcessor) ProcessUpload(ctx context.Context, event UploadedEvent) error {
	raw, err := fp.storage.Get(ctx, event.ObjectKey)
	if err != nil {
		return fmt.Errorf("fetching raw object %s: %w", event.ObjectKey, err)
	}
	defer raw.Close()

	processedKey := "processed/" + strings.TrimPrefix(event.ObjectKey, "raw/") + ".gz"
	dst, err := fp.storage.Create(ctx, processedKey)
	if err != nil {
		return fmt.Errorf("creating processed object: %w", err)
	}

	checksum, err := compressAndHash(raw, dst)
	if err != nil {
		dst.Close()
		fp.storage.Remove(ctx, processedKey)
		return fmt.Errorf("processing %s: %w", event.ObjectKey, err)
	}
	if err := dst.Close(); err != nil {
		fp.storage.Remove(ctx, processedKey)
		return fmt.Errorf("closing processed object: %w", err)
	}

	pe := ProcessedEvent{
		FileId:       event.FileID,
		Checksum:     checksum,
		ProcessedKey: processedKey,
	}
	if err := fp.publisher.PublishFileProcessed(ctx, pe); err != nil {
		fp.storage.Remove(ctx, processedKey)
		return fmt.Errorf("publishing processed event: %w", err)
	}

	if err := fp.storage.Remove(ctx, event.ObjectKey); err != nil {
		fmt.Printf("warning: failed to remove raw object %s: %v\n", event.ObjectKey, err)
	}

	return nil
}

func compressAndHash(r io.Reader, w io.Writer) (string, error) {
	hasher := sha256.New()
	gz := gzip.NewWriter(w)
	if _, err := io.Copy(io.MultiWriter(hasher, gz), r); err != nil {
		return "", err
	}

	if err := gz.Close(); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
