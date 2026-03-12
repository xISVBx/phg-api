package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type AdjustmentItem struct {
	BaseEntity
	AdjustmentID     uuid.UUID `gorm:"type:uuid;index;not null"`
	SaleItemID       uuid.UUID `gorm:"type:uuid;index;not null"`
	QuantityReturned int
	AmountImpact     float64
}
