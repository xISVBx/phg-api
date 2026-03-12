package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type CashMovementRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.CashMovement, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.CashMovement, error)
	Create(ctx context.Context, item *entities.CashMovement) error
	Update(ctx context.Context, item *entities.CashMovement) error
}
