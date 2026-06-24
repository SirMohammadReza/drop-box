package file

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Provider interface {
	NewFile(ctx context.Context, fi NewFileInputs) (*File, error)
}

type fileRepository interface {
	Create(c context.Context, f *File) error
	GetFilesByFolder(c context.Context, folderID *bson.ObjectID) ([]File, error)
	Move(c context.Context, fileID *bson.ObjectID, newFolderID *bson.ObjectID) error
	Copy(c context.Context, f File, destinationFolderID *bson.ObjectID) error
}

type FileService struct {
	fileRepository fileRepository
}

func NewFileService(fr fileRepository) *FileService {
	return &FileService{
		fileRepository: fr,
	}
}

func (f *FileService) NewFile(ctx context.Context, fi NewFileInputs) (*File, error) {
	nf := File{
		Status:    FileStatusUploading,
		FolderID:  fi.FolderID,
		ObjectKey: fi.ObjectKey,
		Name:      fi.Name,
		Size:      fi.Size,
		MimeType:  fi.MimeType,
		UserID:    fi.UserID,
	}

	err := f.fileRepository.Create(ctx, &nf)
	if err != nil {
		return nil, err
	}
	return &nf, nil
}
