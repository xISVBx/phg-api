package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type Adjustment struct {
	BaseEntity
	SaleID          uuid.UUID `gorm:"type:uuid;index;not null"`
	Type            string    `gorm:"index;not null"`
	Reason          string
	PolicyResult    string
	AmountImpact    float64
	CreatedByUserID uuid.UUID `gorm:"type:uuid;index;not null"`
}
