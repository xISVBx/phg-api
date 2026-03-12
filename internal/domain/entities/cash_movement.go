package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type CashMovement struct {
	BaseEntity
	SessionID         *uuid.UUID `gorm:"type:uuid;index"`
	Type              string     `gorm:"index;not null"`
	CategoryID        uuid.UUID  `gorm:"type:uuid;index;not null"`
	Method            string     `gorm:"not null"`
	Amount            float64
	Reference         string
	RelatedEntityType string
	RelatedEntityID   string
	Notes             string
	CreatedByUserID   uuid.UUID `gorm:"type:uuid;index;not null"`
}
