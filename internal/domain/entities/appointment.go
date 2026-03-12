package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type Appointment struct {
	BaseEntity
	CustomerID      uuid.UUID  `gorm:"type:uuid;index;not null"`
	SaleID          *uuid.UUID `gorm:"type:uuid;index"`
	ProductID       uuid.UUID  `gorm:"type:uuid;index;not null"`
	StartsAtUtc     time.Time
	EndsAtUtc       *time.Time
	Status          string `gorm:"index;not null"`
	Notes           string
	CreatedByUserID uuid.UUID `gorm:"type:uuid;index;not null"`
}
