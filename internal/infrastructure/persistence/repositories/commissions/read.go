package commissions

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) ListByWorker(ctx context.Context, workerID uuid.UUID) ([]entities.CommissionEntry, error) {
	var out []entities.CommissionEntry
	err := r.db.WithContext(ctx).Order("earned_at_utc asc").Find(&out, "worker_id = ?", workerID).Error
	return out, err
}

func (r *Repository) ListPendingByWorker(ctx context.Context, workerID uuid.UUID) ([]entities.CommissionEntry, error) {
	var out []entities.CommissionEntry
	err := r.db.WithContext(ctx).Order("earned_at_utc asc").Find(&out, "worker_id = ? AND status IN ?", workerID, []string{"Earned", "PartiallyPaid"}).Error
	return out, err
}
