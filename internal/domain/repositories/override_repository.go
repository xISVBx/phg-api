package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type OverrideRepository interface {
	ListByUser(ctx context.Context, userID uuid.UUID) ([]entities.UserPermissionOverride, error)
	ReplaceByUser(ctx context.Context, userID uuid.UUID, items []entities.UserPermissionOverride) error
	Create(ctx context.Context, item *entities.UserPermissionOverride) error
	Delete(ctx context.Context, userID, overrideID uuid.UUID) error
}
