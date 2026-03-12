package security

import "github.com/google/uuid"

type RolePermissionDTO struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}
