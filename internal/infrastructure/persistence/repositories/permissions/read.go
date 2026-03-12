package permissions

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.Permission, int64, error) {
	var out []entities.Permission
	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"code":         "code",
		"name":         "name",
		"description":  "description",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.Permission{}, &out, opts, []string{"code", "name", "description"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Permission, error) {
	var out entities.Permission
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
