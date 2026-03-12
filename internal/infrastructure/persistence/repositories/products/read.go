package products

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.Product, int64, error) {
	var out []entities.Product
	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"name":         "name",
		"type":         "type",
		"basePrice":    "base_price",
		"cost":         "cost",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.Product{}, &out, opts, []string{"name", "type", "notes"}, allowedSorts)
	return out, total, err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	var out entities.Product
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
