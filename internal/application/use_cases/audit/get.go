package audit

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) Get(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error) {
	return u.uow.Repositories().AuditLogs().GetByID(ctx, id)
}
