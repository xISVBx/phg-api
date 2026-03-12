package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type MenuRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.Menu, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Menu, error)
	Create(ctx context.Context, item *entities.Menu) error
	Update(ctx context.Context, item *entities.Menu) error
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
}
