package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type CommissionRepository interface {
	ListByWorker(ctx context.Context, workerID uuid.UUID) ([]entities.CommissionEntry, error)
	ListPendingByWorker(ctx context.Context, workerID uuid.UUID) ([]entities.CommissionEntry, error)
	CreateMany(ctx context.Context, items []entities.CommissionEntry) error
	Update(ctx context.Context, item *entities.CommissionEntry) error
}
