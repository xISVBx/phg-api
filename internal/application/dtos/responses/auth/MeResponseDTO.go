package auth

import (
	"time"

	"github.com/google/uuid"

	"photogallery/api_go/internal/application/dtos"
)

type MeResponseDTO struct {
	ID           uuid.UUID                      `json:"id"`
	Username     string                         `json:"username"`
	FullName     string                         `json:"fullName"`
	Phone        string                         `json:"phone"`
	Email        string                         `json:"email"`
	IsActive     bool                           `json:"isActive"`
	CreatedAtUtc time.Time                      `json:"createdAtUtc"`
	UpdatedAtUtc time.Time                      `json:"updatedAtUtc"`
	Roles        []MeRoleDTO                    `json:"roles"`
	Permissions  []dtos.EffectivePermissionNode `json:"permissions"`
}
