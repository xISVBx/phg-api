package security

import (
	"context"
	"time"

	"github.com/google/uuid"

	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
)

func (u *UseCase) ListMenus(ctx context.Context, o drepo.QueryOptions) ([]entities.Menu, int64, error) {
	return u.uow.Repositories().Menus().List(ctx, o)
}
func (u *UseCase) GetMenu(ctx context.Context, id uuid.UUID) (*entities.Menu, error) {
	return u.uow.Repositories().Menus().GetByID(ctx, id)
}
func (u *UseCase) CreateMenu(ctx context.Context, actor uuid.UUID, item *entities.Menu) error {
	item.CreatedAtUtc = time.Now().UTC()
	err := u.uow.Repositories().Menus().Create(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Menu", item.ID.String(), "CREATE", item)
	}
	return err
}
func (u *UseCase) UpdateMenu(ctx context.Context, actor uuid.UUID, item *entities.Menu) error {
	item.UpdatedAtUtc = time.Now().UTC()
	err := u.uow.Repositories().Menus().Update(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Menu", item.ID.String(), "UPDATE", item)
	}
	return err
}
func (u *UseCase) SetMenuActive(ctx context.Context, actor, id uuid.UUID, active bool) error {
	err := u.uow.Repositories().Menus().SetActive(ctx, id, active)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Menu", id.String(), map[bool]string{true: "ACTIVATE", false: "DEACTIVATE"}[active], nil)
	}
	return err
}
