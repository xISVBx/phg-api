package workerpayments

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) Create(ctx context.Context, item *entities.WorkerPayment) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) CreateAllocations(ctx context.Context, items []entities.WorkerPaymentAllocation) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}
