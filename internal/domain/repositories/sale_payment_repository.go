package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type SalePaymentRepository interface {
	ListBySale(ctx context.Context, saleID uuid.UUID) ([]entities.SalePayment, error)
	Create(ctx context.Context, payment *entities.SalePayment) error
	SumBySale(ctx context.Context, saleID uuid.UUID) (float64, error)
}
