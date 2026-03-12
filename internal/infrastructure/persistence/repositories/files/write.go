package files

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) Create(ctx context.Context, item *entities.File) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.File{}, "id = ?", id).Error
}
