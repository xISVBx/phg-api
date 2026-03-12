package workorder

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) Get(ctx context.Context, id uuid.UUID) (*entities.WorkOrder, error) {
	return u.uow.Repositories().WorkOrders().GetByID(ctx, id)
}
