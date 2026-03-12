package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type CashSession struct {
	BaseEntity
	OpenedByUserID       uuid.UUID `gorm:"type:uuid;index;not null"`
	OpenedAtUtc          time.Time
	OpeningAmount        float64
	ClosedByUserID       *uuid.UUID `gorm:"type:uuid;index"`
	ClosedAtUtc          *time.Time
	ClosingCountedAmount *float64
	Difference           *float64
	Status               string
}
