package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type WorkOrderItem struct {
	BaseEntity
	WorkOrderID uuid.UUID `gorm:"type:uuid;index;not null"`
	SaleItemID  uuid.UUID `gorm:"type:uuid;index;not null"`
	Status      string
	DueDateUtc  *time.Time
	Notes       string
}
