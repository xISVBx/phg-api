package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type WorkOrderRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.WorkOrder, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.WorkOrder, error)
	Create(ctx context.Context, item *entities.WorkOrder) error
	Update(ctx context.Context, item *entities.WorkOrder) error
}
