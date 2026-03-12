package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type CategoryRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.Category, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Category, error)
	Create(ctx context.Context, item *entities.Category) error
	Update(ctx context.Context, item *entities.Category) error
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
}
