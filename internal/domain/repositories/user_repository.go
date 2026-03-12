package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type UserRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.User, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*entities.User, error)
	Create(ctx context.Context, item *entities.User) error
	Update(ctx context.Context, item *entities.User) error
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
	SetPassword(ctx context.Context, id uuid.UUID, hash string) error
}
