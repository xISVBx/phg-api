package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type WorkerPayment struct {
	BaseEntity
	WorkerID        uuid.UUID `gorm:"type:uuid;index;not null"`
	Type            string    `gorm:"index;not null"`
	Method          string    `gorm:"not null"`
	Amount          float64
	Notes           string
	PaidAtUtc       time.Time
	CreatedByUserID uuid.UUID  `gorm:"type:uuid;index;not null"`
	CashMovementID  *uuid.UUID `gorm:"type:uuid;index"`
}
