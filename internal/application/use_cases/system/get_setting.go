package system

import (
	"context"

	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) GetSetting(ctx context.Context, key string) (*entities.AppSetting, error) {
	return u.uow.Repositories().Settings().Get(ctx, key)
}
