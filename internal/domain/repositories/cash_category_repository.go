package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type CashCategoryRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.CashCategory, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.CashCategory, error)
	Create(ctx context.Context, item *entities.CashCategory) error
	Update(ctx context.Context, item *entities.CashCategory) error
}
