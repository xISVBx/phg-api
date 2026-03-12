package sales

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) Get(ctx context.Context, id uuid.UUID) (*entities.Sale, []entities.SaleItem, []entities.SalePayment, error) {
	sale, err := u.uow.Repositories().Sales().GetByID(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	items, _ := u.uow.Repositories().SaleItems().ListBySale(ctx, id)
	payments, _ := u.uow.Repositories().SalePayments().ListBySale(ctx, id)
	return sale, items, payments, nil
}
