package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type AppointmentRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.Appointment, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Appointment, error)
	Create(ctx context.Context, item *entities.Appointment) error
	Update(ctx context.Context, item *entities.Appointment) error
}
