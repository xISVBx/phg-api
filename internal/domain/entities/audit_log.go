package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type AuditLog struct {
	BaseEntity
	ActorUserID *uuid.UUID `gorm:"type:uuid"`
	EntityType  string     `gorm:"index;not null"`
	EntityID    string     `gorm:"index;not null"`
	Action      string     `gorm:"index;not null"`
	DataJSON    string     `gorm:"type:text"`
	IPAddress   string
}
