package repositories

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
)

type AppSettingRepository interface {
	Get(ctx context.Context, key string) (*entities.AppSetting, error)
	Set(ctx context.Context, item *entities.AppSetting) error
}
