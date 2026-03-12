package cash

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) CreateCategory(ctx context.Context, item *entities.CashCategory) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) UpdateCategory(ctx context.Context, item *entities.CashCategory) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *Repository) CreateMovement(ctx context.Context, item *entities.CashMovement) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) UpdateMovement(ctx context.Context, item *entities.CashMovement) error {
	return r.db.WithContext(ctx).Save(item).Error
}
