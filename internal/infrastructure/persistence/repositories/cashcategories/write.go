package cashcategories

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) Create(ctx context.Context, item *entities.CashCategory) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) Update(ctx context.Context, item *entities.CashCategory) error {
	return r.db.WithContext(ctx).Save(item).Error
}
