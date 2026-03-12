package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type Menu struct {
	BaseEntity
	Code         string `gorm:"uniqueIndex;not null"`
	Name         string `gorm:"not null"`
	DisplayOrder int
	IsActive     bool `gorm:"default:true"`
}
