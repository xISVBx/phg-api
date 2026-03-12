package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type CommissionEntry struct {
	BaseEntity
	WorkerID    uuid.UUID `gorm:"type:uuid;index;not null"`
	SaleItemID  uuid.UUID `gorm:"type:uuid;index;not null"`
	EarnedAtUtc time.Time
	Amount      float64
	PaidAmount  float64
	Status      string `gorm:"index;not null"`
}
