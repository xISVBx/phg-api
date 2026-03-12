package security

import (
	"time"

	"github.com/google/uuid"
)

type UserResponseDTO struct {
	ID           uuid.UUID     `json:"id"`
	Username     string        `json:"username"`
	FullName     string        `json:"fullName"`
	Phone        string        `json:"phone"`
	Email        string        `json:"email"`
	IsActive     bool          `json:"isActive"`
	CreatedAtUtc time.Time     `json:"createdAtUtc"`
	UpdatedAtUtc time.Time     `json:"updatedAtUtc"`
	Roles        []UserRoleDTO `json:"roles,omitempty"`
}
