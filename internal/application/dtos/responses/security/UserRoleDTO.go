package security

import (
	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/enums"
)

type UserRoleDTO struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	IsActive    bool           `json:"isActive"`
	IsPrimary   bool           `json:"isPrimary"`
	RoleType    enums.RoleType `json:"roleType"`
	Menus       []RoleMenuDTO  `json:"menus"`
}
