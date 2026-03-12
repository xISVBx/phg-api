package filelinks

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) Create(ctx context.Context, item *entities.FileLink) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) DeleteByFile(ctx context.Context, fileID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("file_id = ?", fileID).Delete(&entities.FileLink{}).Error
}
