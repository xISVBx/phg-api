package settings

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
)

func (r *Repository) Get(ctx context.Context, key string) (*entities.AppSetting, error) {
	var out entities.AppSetting
	if err := r.db.WithContext(ctx).First(&out, "key = ?", key).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
