package saleitems

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) ListBySale(ctx context.Context, saleID uuid.UUID) ([]entities.SaleItem, error) {
	var out []entities.SaleItem
	err := r.db.WithContext(ctx).Find(&out, "sale_id = ?", saleID).Error
	return out, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.SaleItem, error) {
	var out entities.SaleItem
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
