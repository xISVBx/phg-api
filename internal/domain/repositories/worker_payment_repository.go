package repositories

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
)

type WorkerPaymentRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.WorkerPayment, int64, error)
	Create(ctx context.Context, item *entities.WorkerPayment) error
	CreateAllocations(ctx context.Context, items []entities.WorkerPaymentAllocation) error
}
