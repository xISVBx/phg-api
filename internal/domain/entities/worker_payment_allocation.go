package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type WorkerPaymentAllocation struct {
	BaseEntity
	WorkerPaymentID   uuid.UUID `gorm:"type:uuid;index;not null"`
	CommissionEntryID uuid.UUID `gorm:"type:uuid;index;not null"`
	AmountApplied     float64
}
