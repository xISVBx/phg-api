package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type BaseEntity struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAtUtc time.Time
	UpdatedAtUtc time.Time
}
