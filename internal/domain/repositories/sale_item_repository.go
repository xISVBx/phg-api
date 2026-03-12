package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type SaleItemRepository interface {
	ListBySale(ctx context.Context, saleID uuid.UUID) ([]entities.SaleItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.SaleItem, error)
	CreateMany(ctx context.Context, items []entities.SaleItem) error
}
