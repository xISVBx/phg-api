package cash

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) GetMovement(ctx context.Context, id uuid.UUID) (*entities.CashMovement, error) {
	return u.uow.Repositories().CashMovements().GetByID(ctx, id)
}
