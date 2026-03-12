package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type User struct {
	BaseEntity
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	FullName     string `gorm:"not null"`
	Phone        string
	Email        string
	IsActive     bool `gorm:"default:true"`
}
