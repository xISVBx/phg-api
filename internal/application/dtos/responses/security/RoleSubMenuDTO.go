package security

import "github.com/google/uuid"

type RoleSubMenuDTO struct {
	ID          uuid.UUID           `json:"id"`
	Code        string              `json:"code"`
	Name        string              `json:"name"`
	Permissions []RolePermissionDTO `json:"permissions"`
}
