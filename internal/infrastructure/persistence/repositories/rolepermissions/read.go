package rolepermissions

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) ListRolePermissions(ctx context.Context, roleID uuid.UUID) ([]entities.RoleSubMenuPermission, error) {
	var out []entities.RoleSubMenuPermission
	err := r.db.WithContext(ctx).Find(&out, "role_id = ?", roleID).Error
	return out, err
}
