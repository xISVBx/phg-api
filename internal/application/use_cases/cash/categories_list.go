package cash

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
)

func (u *UseCase) ListCategories(ctx context.Context, o drepo.QueryOptions) ([]entities.CashCategory, int64, error) {
	return u.uow.Repositories().CashCategories().List(ctx, o)
}
