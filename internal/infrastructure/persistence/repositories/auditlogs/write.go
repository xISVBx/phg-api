package auditlogs

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) Create(ctx context.Context, item *entities.AuditLog) error {
	return r.db.WithContext(ctx).Create(item).Error
}
