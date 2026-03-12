package catalog

import (
	"context"
	"time"

	"github.com/google/uuid"
	catalogreq "photogallery/api_go/internal/application/dtos/request/catalog"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) UpdateProduct(ctx context.Context, actor, id uuid.UUID, in catalogreq.UpdateProductRequestDTO) (*entities.Product, error) {
	item, err := u.uow.Repositories().Products().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	catID, err := uuid.Parse(in.CategoryID)
	if err != nil {
		return nil, err
	}
	item.CategoryID, item.Name, item.Type, item.BasePrice, item.Cost, item.CommissionType, item.CommissionValue, item.RequiresDelivery, item.DefaultLeadDays, item.Notes = catID, in.Name, in.Type, in.BasePrice, in.Cost, in.CommissionType, in.CommissionValue, in.RequiresDelivery, in.DefaultLeadDays, in.Notes
	item.UpdatedAtUtc = time.Now().UTC()
	err = u.uow.Repositories().Products().Update(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Product", item.ID.String(), "UPDATE", item)
	}
	return item, err
}
