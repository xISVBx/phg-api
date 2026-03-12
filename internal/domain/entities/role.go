package entities

import (
	"time"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/enums"
)

var _ = time.Now
var _ = uuid.Nil

type Role struct {
	BaseEntity
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	IsActive    bool           `gorm:"default:true"`
	RoleType    enums.RoleType `gorm:"type:varchar(32);not null;default:SYSTEM" json:"roleType"`
}
