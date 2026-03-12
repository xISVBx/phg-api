package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type Category struct {
	BaseEntity
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	IsActive    bool `gorm:"default:true"`
}
