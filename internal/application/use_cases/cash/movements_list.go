package cash

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
)

func (u *UseCase) ListMovements(ctx context.Context, o drepo.QueryOptions) ([]entities.CashMovement, int64, error) {
	return u.uow.Repositories().CashMovements().List(ctx, o)
}
