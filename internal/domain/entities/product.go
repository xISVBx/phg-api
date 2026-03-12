package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type Product struct {
	BaseEntity
	CategoryID       uuid.UUID `gorm:"type:uuid;index;not null"`
	Name             string    `gorm:"not null"`
	Type             string    `gorm:"not null"`
	BasePrice        float64
	Cost             float64
	CommissionType   string
	CommissionValue  float64
	RequiresDelivery bool
	DefaultLeadDays  *int
	IsActive         bool `gorm:"default:true"`
	Notes            string
}
