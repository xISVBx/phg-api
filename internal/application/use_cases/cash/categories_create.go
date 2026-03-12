package cash

import (
	"context"

	"github.com/google/uuid"
	cashreq "photogallery/api_go/internal/application/dtos/request/cash"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) CreateCategory(ctx context.Context, actor uuid.UUID, in cashreq.CreateCashCategoryRequestDTO) (*entities.CashCategory, error) {
	item := &entities.CashCategory{Name: in.Name, Type: in.Type, IsActive: true}
	err := u.uow.Repositories().CashCategories().Create(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "CashCategory", item.ID.String(), "CREATE", item)
	}
	return item, err
}
