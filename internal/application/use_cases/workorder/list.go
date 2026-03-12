package workorder

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
)

func (u *UseCase) List(ctx context.Context, o drepo.QueryOptions) ([]entities.WorkOrder, int64, error) {
	return u.uow.Repositories().WorkOrders().List(ctx, o)
}
