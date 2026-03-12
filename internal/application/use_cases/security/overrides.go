package security

import (
	"context"

	"github.com/google/uuid"

	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) ListOverrides(ctx context.Context, userID uuid.UUID) ([]entities.UserPermissionOverride, error) {
	return u.uow.Repositories().Overrides().ListByUser(ctx, userID)
}
func (u *UseCase) ReplaceOverrides(ctx context.Context, actor, userID uuid.UUID, items []entities.UserPermissionOverride) error {
	err := u.uow.Repositories().Overrides().ReplaceByUser(ctx, userID, items)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "UserOverride", userID.String(), "REPLACE", items)
	}
	return err
}
func (u *UseCase) CreateOverride(ctx context.Context, actor uuid.UUID, item *entities.UserPermissionOverride) error {
	err := u.uow.Repositories().Overrides().Create(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "UserOverride", item.UserID.String(), "CREATE", item)
	}
	return err
}
func (u *UseCase) DeleteOverride(ctx context.Context, actor, userID, overrideID uuid.UUID) error {
	err := u.uow.Repositories().Overrides().Delete(ctx, userID, overrideID)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "UserOverride", overrideID.String(), "DELETE", nil)
	}
	return err
}
