package repositories

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

type AuditLogRepository interface {
	List(ctx context.Context, opts QueryOptions) ([]entities.AuditLog, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error)
	Create(ctx context.Context, item *entities.AuditLog) error
}
