package auth

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"photogallery/api_go/internal/application/dtos"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/enums"
)

func (u *UseCase) MyPermissions(ctx context.Context, userID uuid.UUID) ([]dtos.EffectivePermissionNode, error) {
	repos := u.uow.Repositories()
	urs, err := repos.UserRoles().ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(urs) == 0 {
		return []dtos.EffectivePermissionNode{}, nil
	}
	roleID := urs[0].RoleID
	base, err := repos.RolePermissions().ListRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}
	overrides, _ := repos.Overrides().ListByUser(ctx, userID)
	set := map[string]bool{}
	for _, p := range base {
		set[p.SubMenuID.String()+":"+p.PermissionID.String()] = true
	}
	for _, o := range overrides {
		key := o.SubMenuID.String() + ":" + o.PermissionID.String()
		if strings.EqualFold(o.Mode, string(enums.OverrideGrant)) {
			set[key] = true
		} else {
			delete(set, key)
		}
	}
	subMenus, _, _ := repos.SubMenus().List(ctx, common.QueryOpts(1, 2000, "", "", ""))
	menus, _, _ := repos.Menus().List(ctx, common.QueryOpts(1, 200, "", "", ""))
	perms, _, _ := repos.Permissions().List(ctx, common.QueryOpts(1, 2000, "", "", ""))
	permName := map[uuid.UUID]string{}
	for _, p := range perms {
		permName[p.ID] = p.Code
	}
	byMenu := map[uuid.UUID][]dtos.EffectiveSubMenu{}
	for _, sm := range subMenus {
		item := dtos.EffectiveSubMenu{SubMenuCode: sm.Code, SubMenuName: sm.Name, Permissions: []string{}}
		for k, ok := range set {
			if !ok || !strings.HasPrefix(k, sm.ID.String()+":") {
				continue
			}
			parts := strings.Split(k, ":")
			pid, err := uuid.Parse(parts[1])
			if err == nil {
				item.Permissions = append(item.Permissions, permName[pid])
			}
		}
		if len(item.Permissions) > 0 {
			byMenu[sm.MenuID] = append(byMenu[sm.MenuID], item)
		}
	}
	out := make([]dtos.EffectivePermissionNode, 0)
	for _, m := range menus {
		if len(byMenu[m.ID]) == 0 {
			continue
		}
		out = append(out, dtos.EffectivePermissionNode{MenuCode: m.Code, MenuName: m.Name, SubMenus: byMenu[m.ID]})
	}
	return out, nil
}
