package files

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) ResolveDownloadPath(ctx context.Context, fileID uuid.UUID) (string, *entities.File, error) {
	file, err := u.uow.Repositories().Files().GetByID(ctx, fileID)
	if err != nil {
		return "", nil, err
	}
	path, err := u.storage.OpenAbsPath(file.StorageRelativePath)
	return path, file, err
}
