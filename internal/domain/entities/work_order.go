package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type WorkOrder struct {
	BaseEntity
	SaleID            uuid.UUID `gorm:"type:uuid;index;not null"`
	Status            string    `gorm:"index;not null"`
	DueDateUtc        *time.Time
	ResponsibleUserID *uuid.UUID `gorm:"type:uuid;index"`
	Notes             string
	DeliveredAtUtc    *time.Time
}
