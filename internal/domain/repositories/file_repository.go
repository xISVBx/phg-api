package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type FileRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.File, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.File, error)
	Create(ctx context.Context, item *entities.File) error
	Delete(ctx context.Context, id uuid.UUID) error
}
