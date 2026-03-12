package overrides

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) ReplaceByUser(ctx context.Context, userID uuid.UUID, items []entities.UserPermissionOverride) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entities.UserPermissionOverride{}).Error; err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *Repository) Create(ctx context.Context, item *entities.UserPermissionOverride) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) Delete(ctx context.Context, userID, overrideID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", overrideID, userID).Delete(&entities.UserPermissionOverride{}).Error
}
