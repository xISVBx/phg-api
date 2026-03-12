package appointment

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) Get(ctx context.Context, id uuid.UUID) (*entities.Appointment, error) {
	return u.uow.Repositories().Appointments().GetByID(ctx, id)
}
