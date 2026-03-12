package filelinks

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) ListByFile(ctx context.Context, fileID uuid.UUID) ([]entities.FileLink, error) {
	var out []entities.FileLink
	err := r.db.WithContext(ctx).Find(&out, "file_id = ?", fileID).Error
	return out, err
}
