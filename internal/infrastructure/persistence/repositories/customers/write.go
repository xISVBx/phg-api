package customers

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) Create(ctx context.Context, item *entities.Customer) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) Update(ctx context.Context, item *entities.Customer) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *Repository) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	return r.db.WithContext(ctx).Model(&entities.Customer{}).Where("id = ?", id).Update("is_active", active).Error
}
