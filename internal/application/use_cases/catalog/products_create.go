package catalog

import (
	"context"

	"github.com/google/uuid"
	catalogreq "photogallery/api_go/internal/application/dtos/request/catalog"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) CreateProduct(ctx context.Context, actor uuid.UUID, in catalogreq.CreateProductRequestDTO) (*entities.Product, error) {
	catID, err := uuid.Parse(in.CategoryID)
	if err != nil {
		return nil, err
	}
	item := &entities.Product{CategoryID: catID, Name: in.Name, Type: in.Type, BasePrice: in.BasePrice, Cost: in.Cost, CommissionType: in.CommissionType, CommissionValue: in.CommissionValue, RequiresDelivery: in.RequiresDelivery, DefaultLeadDays: in.DefaultLeadDays, Notes: in.Notes, IsActive: true}
	err = u.uow.Repositories().Products().Create(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Product", item.ID.String(), "CREATE", item)
	}
	return item, err
}
