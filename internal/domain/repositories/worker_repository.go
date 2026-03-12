package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type WorkerRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.Worker, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Worker, error)
	Create(ctx context.Context, item *entities.Worker) error
	Update(ctx context.Context, item *entities.Worker) error
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
}
