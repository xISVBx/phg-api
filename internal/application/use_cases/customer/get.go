package customer

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) Get(ctx context.Context, id uuid.UUID) (*entities.Customer, error) {
	return u.uow.Repositories().Customers().GetByID(ctx, id)
}
