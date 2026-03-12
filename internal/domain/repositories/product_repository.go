package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type ProductRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.Product, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	Create(ctx context.Context, item *entities.Product) error
	Update(ctx context.Context, item *entities.Product) error
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
}
