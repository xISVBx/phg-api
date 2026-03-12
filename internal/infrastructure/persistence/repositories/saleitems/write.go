package saleitems

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) CreateMany(ctx context.Context, items []entities.SaleItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}
