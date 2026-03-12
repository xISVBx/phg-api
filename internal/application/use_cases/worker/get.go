package worker

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) Get(ctx context.Context, id uuid.UUID) (*entities.Worker, error) {
	return u.uow.Repositories().Workers().GetByID(ctx, id)
}
