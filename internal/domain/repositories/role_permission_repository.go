package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type RolePermissionRepository interface {
	ListRolePermissions(ctx context.Context, roleID uuid.UUID) ([]entities.RoleSubMenuPermission, error)
	ReplaceRolePermissions(ctx context.Context, roleID uuid.UUID, items []entities.RoleSubMenuPermission) error
}
