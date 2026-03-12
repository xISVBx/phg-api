package cash

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) ListCategories(ctx context.Context, opts drepo.QueryOptions) ([]entities.CashCategory, int64, error) {
	var out []entities.CashCategory
	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"name":         "name",
		"type":         "type",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.CashCategory{}, &out, opts, []string{"name", "type"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*entities.CashCategory, error) {
	var out entities.CashCategory
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *Repository) ListMovements(ctx context.Context, opts drepo.QueryOptions) ([]entities.CashMovement, int64, error) {
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

func (r *Repository) GetMovementByID(ctx context.Context, id uuid.UUID) (*entities.CashMovement, error) {
	var out entities.CashMovement
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
