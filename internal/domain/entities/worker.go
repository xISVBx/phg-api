package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type Worker struct {
	BaseEntity
	FullName     string `gorm:"not null"`
	Phone        string
	Email        string
	IsActive     bool `gorm:"default:true"`
	FixedSalary  float64
	SalaryPeriod string
	Notes        string
}
