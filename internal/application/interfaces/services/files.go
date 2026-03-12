package services

import (
	"context"
	"mime/multipart"
)

type StoredFile struct {
	OriginalName        string
	ContentType         string
	SizeBytes           int64
	StorageRelativePath string
	StorageBaseKey      string
}

type IFileStorageService interface {
	Store(ctx context.Context, file *multipart.FileHeader, relativeDir string) (*StoredFile, error)
	Delete(ctx context.Context, relativePath string) error
	OpenAbsPath(relativePath string) (string, error)
}
