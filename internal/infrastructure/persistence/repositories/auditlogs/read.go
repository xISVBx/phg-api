package auditlogs

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.AuditLog, int64, error) {
	var out []entities.AuditLog
	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"entityType":   "entity_type",
		"entityId":     "entity_id",
		"action":       "action",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.AuditLog{}, &out, opts, []string{"entity_type", "entity_id", "action"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error) {
	var out entities.AuditLog
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
