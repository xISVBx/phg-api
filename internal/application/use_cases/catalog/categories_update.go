package catalog

import (
	"context"
	"time"

	"github.com/google/uuid"
	catalogreq "photogallery/api_go/internal/application/dtos/request/catalog"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) UpdateCategory(ctx context.Context, actor, id uuid.UUID, in catalogreq.UpdateCategoryRequestDTO) (*entities.Category, error) {
	item, err := u.uow.Repositories().Categories().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Name, item.Description = in.Name, in.Description
	item.UpdatedAtUtc = time.Now().UTC()
	err = u.uow.Repositories().Categories().Update(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Category", item.ID.String(), "UPDATE", item)
	}
	return item, err
}
