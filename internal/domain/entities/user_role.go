package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type UserRole struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoleID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	IsPrimary bool
}
