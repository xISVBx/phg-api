package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type FileLink struct {
	BaseEntity
	FileID     uuid.UUID  `gorm:"type:uuid;index;not null"`
	EntityType string     `gorm:"index;not null"`
	EntityID   uuid.UUID  `gorm:"type:uuid;index;not null"`
	CustomerID *uuid.UUID `gorm:"type:uuid;index"`
	Notes      string
}
