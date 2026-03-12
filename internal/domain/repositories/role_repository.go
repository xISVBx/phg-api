package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type RoleRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.Role, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Role, error)
	Create(ctx context.Context, item *entities.Role) error
	Update(ctx context.Context, item *entities.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
}
