package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type FileLinkRepository interface {
	ListByFile(ctx context.Context, fileID uuid.UUID) ([]entities.FileLink, error)
	Create(ctx context.Context, item *entities.FileLink) error
	DeleteByFile(ctx context.Context, fileID uuid.UUID) error
}
