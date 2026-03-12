package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type Sale struct {
	BaseEntity
	CustomerID              *uuid.UUID `gorm:"type:uuid;index"`
	SellerUserID            uuid.UUID  `gorm:"type:uuid;index;not null"`
	Status                  string     `gorm:"index;not null"`
	NotifyOptIn             bool
	Subtotal                float64
	DiscountTotal           float64
	Total                   float64
	TotalCostSnapshot       float64
	TotalCommissionSnapshot float64
	ConfirmedAtUtc          *time.Time
}
