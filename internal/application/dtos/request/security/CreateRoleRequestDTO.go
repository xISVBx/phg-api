package security

import "photogallery/api_go/internal/domain/enums"

type CreateRoleRequestDTO struct {
	Name        string                     `json:"name" binding:"required"`
	Description string                     `json:"description"`
	IsActive    *bool                      `json:"isActive"`
	RoleType    enums.RoleType             `json:"roleType"`
	Permissions []RolePermissionSetItemDTO `json:"permissions" binding:"required"`
}
