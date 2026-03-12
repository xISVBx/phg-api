package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type SaleRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.Sale, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Sale, error)
	Create(ctx context.Context, sale *entities.Sale) error
	Update(ctx context.Context, sale *entities.Sale) error
}
