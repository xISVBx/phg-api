package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type CashCategory struct {
	BaseEntity
	Name     string `gorm:"uniqueIndex;not null"`
	Type     string `gorm:"not null"`
	IsActive bool   `gorm:"default:true"`
}
