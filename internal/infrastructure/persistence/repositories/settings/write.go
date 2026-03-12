package settings

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) Set(ctx context.Context, item *entities.AppSetting) error {
	return r.db.WithContext(ctx).Save(item).Error
}
