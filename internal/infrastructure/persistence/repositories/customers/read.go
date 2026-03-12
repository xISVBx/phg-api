package customers

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
)

func normalizePage(opts drepo.QueryOptions) (int, int) {
	page := opts.Page
	if page <= 0 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 25
	}
	return page, pageSize
}

func (r *Repository) List(ctx context.Context, opts drepo.QueryOptions) ([]entities.Customer, int64, error) {
	page, pageSize := normalizePage(opts)
	tx := r.db.WithContext(ctx).Model(&entities.Customer{})
	if q := strings.TrimSpace(opts.Q); q != "" {
		p := "%" + strings.ToLower(q) + "%"
		tx = tx.Where("LOWER(full_name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(customer_code) LIKE ?", p, p, p)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	order := "created_at_utc desc"
	if opts.Sort != "" {
		dir := strings.ToLower(opts.Dir)
		if dir != "asc" {
			dir = "desc"
		}
		order = opts.Sort + " " + dir
	}
	var out []entities.Customer
	if err := tx.Order(order).Limit(pageSize).Offset((page - 1) * pageSize).Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Customer, error) {
	var out entities.Customer
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
