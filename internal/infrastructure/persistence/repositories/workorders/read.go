package workorders

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.WorkOrder, int64, error) {
	var out []entities.WorkOrder
	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"status":       "status",
		"dueDateUtc":   "due_date_utc",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.WorkOrder{}, &out, opts, []string{"status", "notes"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.WorkOrder, error) {
	var out entities.WorkOrder
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
