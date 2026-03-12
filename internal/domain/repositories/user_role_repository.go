package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type UserRoleRepository interface {
	ListByUser(ctx context.Context, userID uuid.UUID) ([]entities.UserRole, error)
	ReplaceByUser(ctx context.Context, userID uuid.UUID, items []entities.UserRole) error
	SetPrimaryRole(ctx context.Context, userID, roleID uuid.UUID) error
}
