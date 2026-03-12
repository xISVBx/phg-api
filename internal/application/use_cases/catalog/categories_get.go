package catalog

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) GetCategory(ctx context.Context, id uuid.UUID) (*entities.Category, error) {
	return u.uow.Repositories().Categories().GetByID(ctx, id)
}
