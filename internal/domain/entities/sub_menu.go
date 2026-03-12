package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type SubMenu struct {
	BaseEntity
	MenuID       uuid.UUID `gorm:"type:uuid;index;not null"`
	Code         string    `gorm:"uniqueIndex;not null"`
	Name         string    `gorm:"not null"`
	Route        string
	DisplayOrder int
	IsActive     bool `gorm:"default:true"`
}
