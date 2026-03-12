package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type SaleItem struct {
	BaseEntity
	SaleID                   uuid.UUID `gorm:"type:uuid;index;not null"`
	ProductID                uuid.UUID `gorm:"type:uuid;index;not null"`
	Quantity                 int       `gorm:"not null"`
	UnitPriceSnapshot        float64
	UnitCostSnapshot         float64
	CommissionTypeSnapshot   string
	CommissionValueSnapshot  float64
	DiscountSnapshot         float64
	DiscountReason           string
	Notes                    string
	RequiresDeliverySnapshot bool
	LeadDaysSnapshot         *int
}
