package userroles

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) ReplaceByUser(ctx context.Context, userID uuid.UUID, items []entities.UserRole) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entities.UserRole{}).Error; err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *Repository) SetPrimaryRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return r.ReplaceByUser(ctx, userID, []entities.UserRole{{UserID: userID, RoleID: roleID, IsPrimary: true}})
}
