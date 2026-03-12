package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type SalePayment struct {
	BaseEntity
	SaleID          uuid.UUID `gorm:"type:uuid;index;not null"`
	Method          string    `gorm:"not null"`
	Amount          float64
	Reference       string
	PaidAtUtc       time.Time
	CreatedByUserID uuid.UUID `gorm:"type:uuid;index;not null"`
}
