package salepayments

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) ListBySale(ctx context.Context, saleID uuid.UUID) ([]entities.SalePayment, error) {
	var out []entities.SalePayment
	err := r.db.WithContext(ctx).Find(&out, "sale_id = ?", saleID).Error
	return out, err
}

func (r *Repository) SumBySale(ctx context.Context, saleID uuid.UUID) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).Model(&entities.SalePayment{}).Select("COALESCE(SUM(amount),0)").Where("sale_id = ?", saleID).Scan(&total).Error
	return total, err
}
