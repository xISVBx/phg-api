package commissions

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) CreateMany(ctx context.Context, items []entities.CommissionEntry) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *Repository) Update(ctx context.Context, item *entities.CommissionEntry) error {
	return r.db.WithContext(ctx).Save(item).Error
}
