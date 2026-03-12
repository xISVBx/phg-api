package entities

import (
	"time"

	"github.com/google/uuid"
)

var _ = time.Now
var _ = uuid.Nil

type File struct {
	BaseEntity
	OriginalName        string `gorm:"not null"`
	ContentType         string `gorm:"not null"`
	SizeBytes           int64
	UploadedByUserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	UploadedAtUtc       time.Time
	StorageRelativePath string `gorm:"not null;index"`
	StorageKind         string `gorm:"not null;index"`
	StorageBaseKey      string
}
