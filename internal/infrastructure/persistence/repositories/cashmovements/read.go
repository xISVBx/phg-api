package cashmovements

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.CashMovement, int64, error) {
	var out []entities.CashMovement
	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"type":         "type",
		"method":       "method",
		"reference":    "reference",
		"amount":       "amount",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.CashMovement{}, &out, opts, []string{"type", "method", "reference"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.CashMovement, error) {
	var out entities.CashMovement
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
