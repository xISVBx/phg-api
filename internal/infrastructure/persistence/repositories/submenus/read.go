package submenus

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.SubMenu, int64, error) {
	var out []entities.SubMenu
	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"code":         "code",
		"name":         "name",
		"route":        "route",
		"displayOrder": "display_order",
		"order":        "display_order",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.SubMenu{}, &out, opts, []string{"code", "name", "route"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.SubMenu, error) {
	var out entities.SubMenu
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
