package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type Customer struct {
	BaseEntity
	FullName     string `gorm:"not null"`
	Phone        string
	Email        string `gorm:"uniqueIndex"`
	CustomerCode string `gorm:"uniqueIndex;not null"`
	Document     string
	Notes        string
	IsActive     bool `gorm:"default:true"`
}
