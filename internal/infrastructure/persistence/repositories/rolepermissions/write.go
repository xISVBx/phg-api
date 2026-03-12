package rolepermissions

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) ReplaceRolePermissions(ctx context.Context, roleID uuid.UUID, items []entities.RoleSubMenuPermission) error {
	if err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Delete(&entities.RoleSubMenuPermission{}).Error; err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}
