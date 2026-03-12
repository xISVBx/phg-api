package system

import (
	"context"
	"time"

	"github.com/google/uuid"

	appif "photogallery/api_go/internal/application/interfaces"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) SetSetting(ctx context.Context, actor uuid.UUID, key, value string) error {
	now := time.Now().UTC()
	return u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
		item := &entities.AppSetting{Key: key, Value: value, UpdatedAtUtc: now, UpdatedByUserID: &actor}
		if err := repos.Settings().Set(ctx, item); err != nil {
			return err
		}
		common.CreateAudit(ctx, repos, &actor, "AppSetting", key, "SET", map[string]any{"value": value})
		return nil
	})
}
