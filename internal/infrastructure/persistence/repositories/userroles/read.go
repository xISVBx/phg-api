package userroles

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) ListByUser(ctx context.Context, userID uuid.UUID) ([]entities.UserRole, error) {
	var out []entities.UserRole
	err := r.db.WithContext(ctx).Find(&out, "user_id = ?", userID).Error
	return out, err
}
