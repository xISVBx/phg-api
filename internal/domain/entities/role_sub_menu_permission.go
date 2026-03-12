package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type RoleSubMenuPermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	SubMenuID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	PermissionID uuid.UUID `gorm:"type:uuid;primaryKey"`
}
