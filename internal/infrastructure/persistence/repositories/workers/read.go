package workers

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.Worker, int64, error) {
	var out []entities.Worker
	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"fullName":     "full_name",
		"email":        "email",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.Worker{}, &out, opts, []string{"full_name", "email"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Worker, error) {
	var out entities.Worker
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
