package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type SubMenuRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.SubMenu, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.SubMenu, error)
	Create(ctx context.Context, item *entities.SubMenu) error
	Update(ctx context.Context, item *entities.SubMenu) error
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
}
