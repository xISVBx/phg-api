package cash

import (
	"context"
	"time"

	"github.com/google/uuid"
	cashreq "photogallery/api_go/internal/application/dtos/request/cash"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) UpdateCategory(ctx context.Context, actor, id uuid.UUID, in cashreq.UpdateCashCategoryRequestDTO) (*entities.CashCategory, error) {
	item, err := u.uow.Repositories().CashCategories().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Name, item.Type = in.Name, in.Type
	item.UpdatedAtUtc = time.Now().UTC()
	err = u.uow.Repositories().CashCategories().Update(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "CashCategory", item.ID.String(), "UPDATE", item)
	}
	return item, err
}
