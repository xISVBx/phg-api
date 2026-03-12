package catalog

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) GetProduct(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	return u.uow.Repositories().Products().GetByID(ctx, id)
}
