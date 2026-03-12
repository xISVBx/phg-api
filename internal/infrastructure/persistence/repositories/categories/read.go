package categories

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.Category, int64, error) {
	var out []entities.Category
	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"name":         "name",
		"description":  "description",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.Category{}, &out, opts, []string{"name", "description"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Category, error) {
	var out entities.Category
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
