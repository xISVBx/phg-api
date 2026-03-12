package overrides

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) ListByUser(ctx context.Context, userID uuid.UUID) ([]entities.UserPermissionOverride, error) {
	var out []entities.UserPermissionOverride
	err := r.db.WithContext(ctx).Find(&out, "user_id = ?", userID).Error
	return out, err
}
