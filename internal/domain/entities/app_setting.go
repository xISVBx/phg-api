package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type AppSetting struct {
	Key             string `gorm:"primaryKey"`
	Value           string `gorm:"type:text"`
	UpdatedAtUtc    time.Time
	UpdatedByUserID *uuid.UUID `gorm:"type:uuid"`
}
