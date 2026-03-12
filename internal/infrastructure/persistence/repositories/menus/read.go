package menus

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.Menu, int64, error) {
	var out []entities.Menu
	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"code":         "code",
		"name":         "name",
		"displayOrder": "display_order",
		"order":        "display_order",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.Menu{}, &out, opts, []string{"code", "name"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Menu, error) {
	var out entities.Menu
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
