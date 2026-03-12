package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type UserPermissionOverride struct {
	BaseEntity
	UserID       uuid.UUID `gorm:"type:uuid;index;not null"`
	SubMenuID    uuid.UUID `gorm:"type:uuid;index;not null"`
	PermissionID uuid.UUID `gorm:"type:uuid;index;not null"`
	Mode         string    `gorm:"not null"`
}
