package catalog

import (
	"context"

	"github.com/google/uuid"
	catalogreq "photogallery/api_go/internal/application/dtos/request/catalog"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) CreateCategory(ctx context.Context, actor uuid.UUID, in catalogreq.CreateCategoryRequestDTO) (*entities.Category, error) {
	item := &entities.Category{Name: in.Name, Description: in.Description, IsActive: true}
	err := u.uow.Repositories().Categories().Create(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Category", item.ID.String(), "CREATE", item)
	}
	return item, err
}
