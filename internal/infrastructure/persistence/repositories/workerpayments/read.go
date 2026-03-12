package workerpayments

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
	repocommon "photogallery/api_go/internal/infrastructure/persistence/repositories/common"
)

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.WorkerPayment, int64, error) {
	var out []entities.WorkerPayment
	allowedSorts := map[string]string{
		"createdAtUtc": "created_at_utc",
		"paidAtUtc":    "paid_at_utc",
		"type":         "type",
		"method":       "method",
		"amount":       "amount",
	}
	total, err := repocommon.ListWithQuery(ctx, r.db, &entities.WorkerPayment{}, &out, opts, []string{"type", "method"}, allowedSorts)
	return out, total, err
}
