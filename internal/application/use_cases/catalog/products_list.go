package catalog

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
)

func (u *UseCase) ListProducts(ctx context.Context, o drepo.QueryOptions) ([]entities.Product, int64, error) {
	return u.uow.Repositories().Products().List(ctx, o)
}
