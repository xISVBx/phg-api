package roles

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) Create(ctx context.Context, item *entities.Role) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) Update(ctx context.Context, item *entities.Role) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Role{}, "id = ?", id).Error
}

func (r *Repository) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	return r.db.WithContext(ctx).Model(&entities.Role{}).Where("id = ?", id).Update("is_active", active).Error
}
