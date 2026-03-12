package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type PermissionRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.Permission, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Permission, error)
	Create(ctx context.Context, item *entities.Permission) error
	Update(ctx context.Context, item *entities.Permission) error
}
