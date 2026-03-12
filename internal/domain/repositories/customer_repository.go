package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type CustomerRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.Customer, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Customer, error)
	Create(ctx context.Context, item *entities.Customer) error
	Update(ctx context.Context, item *entities.Customer) error
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
}
