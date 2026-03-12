package security

import (
	"context"
	"time"

	"github.com/google/uuid"

	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
	drepo "photogallery/api_go/internal/domain/repositories"
)

func (u *UseCase) ListSubMenus(ctx context.Context, o drepo.QueryOptions) ([]entities.SubMenu, int64, error) {
	return u.uow.Repositories().SubMenus().List(ctx, o)
}
func (u *UseCase) GetSubMenu(ctx context.Context, id uuid.UUID) (*entities.SubMenu, error) {
	return u.uow.Repositories().SubMenus().GetByID(ctx, id)
}
func (u *UseCase) CreateSubMenu(ctx context.Context, actor uuid.UUID, item *entities.SubMenu) error {
	item.CreatedAtUtc = time.Now().UTC()
	err := u.uow.Repositories().SubMenus().Create(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "SubMenu", item.ID.String(), "CREATE", item)
	}
	return err
}
func (u *UseCase) UpdateSubMenu(ctx context.Context, actor uuid.UUID, item *entities.SubMenu) error {
	item.UpdatedAtUtc = time.Now().UTC()
	err := u.uow.Repositories().SubMenus().Update(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "SubMenu", item.ID.String(), "UPDATE", item)
	}
	return err
}
func (u *UseCase) SetSubMenuActive(ctx context.Context, actor, id uuid.UUID, active bool) error {
	err := u.uow.Repositories().SubMenus().SetActive(ctx, id, active)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "SubMenu", id.String(), map[bool]string{true: "ACTIVATE", false: "DEACTIVATE"}[active], nil)
	}
	return err
}

func (u *UseCase) ListPermissions(ctx context.Context, o drepo.QueryOptions) ([]entities.Permission, int64, error) {
	return u.uow.Repositories().Permissions().List(ctx, o)
}
func (u *UseCase) GetPermission(ctx context.Context, id uuid.UUID) (*entities.Permission, error) {
	return u.uow.Repositories().Permissions().GetByID(ctx, id)
}
func (u *UseCase) CreatePermission(ctx context.Context, actor uuid.UUID, item *entities.Permission) error {
	item.CreatedAtUtc = time.Now().UTC()
	err := u.uow.Repositories().Permissions().Create(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Permission", item.ID.String(), "CREATE", item)
	}
	return err
}
func (u *UseCase) UpdatePermission(ctx context.Context, actor uuid.UUID, item *entities.Permission) error {
	item.UpdatedAtUtc = time.Now().UTC()
	err := u.uow.Repositories().Permissions().Update(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Permission", item.ID.String(), "UPDATE", item)
	}
	return err
}
