package roles

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/domain/enums"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.Role, int64, error) {
	var out []entities.Role

	db := r.db.
		Where("role_type = ?", enums.RoleTypeSystem)

	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"name":         "name",
		"description":  "description",
		"roleType":     "role_type",
	}
	total, err := repocommon.ListWithQuery(ctx, db, &entities.Role{}, &out, opts, []string{"name", "description"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Role, error) {
	var out entities.Role
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
