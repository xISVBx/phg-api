package catalog

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
)

func (u *UseCase) ListCategories(ctx context.Context, o drepo.QueryOptions) ([]entities.Category, int64, error) {
	return u.uow.Repositories().Categories().List(ctx, o)
}
