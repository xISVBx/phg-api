package security

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	securityreq "photogallery/api_go/internal/application/dtos/request/security"
	securityresp "photogallery/api_go/internal/application/dtos/responses/security"
	appif "photogallery/api_go/internal/application/interfaces"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/domain/enums"
	drepo "photogallery/api_go/internal/domain/repositories"
)

func (u *UseCase) ListRoles(ctx context.Context, o drepo.QueryOptions) ([]securityresp.RoleResponseDTO, int64, error) {
	repos := u.uow.Repositories()
	roles, total, err := repos.Roles().List(ctx, o)
	if err != nil {
		return nil, 0, err
	}

	menus, _, err := repos.Menus().List(ctx, common.QueryOpts(1, 5000, "", "display_order", "asc"))
	if err != nil {
		return nil, 0, err
	}
	subMenus, _, err := repos.SubMenus().List(ctx, common.QueryOpts(1, 10000, "", "display_order", "asc"))
	if err != nil {
		return nil, 0, err
	}
	perms, _, err := repos.Permissions().List(ctx, common.QueryOpts(1, 5000, "", "code", "asc"))
	if err != nil {
		return nil, 0, err
	}
	menuByID, subMenuByID, permByID := buildRoleRefMaps(menus, subMenus, perms)

	out := make([]securityresp.RoleResponseDTO, 0, len(roles))
	for _, role := range roles {
		rolePerms, err := repos.RolePermissions().ListRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, buildRoleResponseDTO(role, rolePerms, menuByID, subMenuByID, permByID))
	}

	return out, total, nil
}
func (u *UseCase) GetRole(ctx context.Context, id uuid.UUID) (*securityresp.RoleResponseDTO, error) {
	repos := u.uow.Repositories()
	role, err := repos.Roles().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	menus, _, err := repos.Menus().List(ctx, common.QueryOpts(1, 5000, "", "display_order", "asc"))
	if err != nil {
		return nil, err
	}
	subMenus, _, err := repos.SubMenus().List(ctx, common.QueryOpts(1, 10000, "", "display_order", "asc"))
	if err != nil {
		return nil, err
	}
	perms, _, err := repos.Permissions().List(ctx, common.QueryOpts(1, 5000, "", "code", "asc"))
	if err != nil {
		return nil, err
	}
	rolePerms, err := repos.RolePermissions().ListRolePermissions(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	menuByID, subMenuByID, permByID := buildRoleRefMaps(menus, subMenus, perms)
	out := buildRoleResponseDTO(*role, rolePerms, menuByID, subMenuByID, permByID)
	return &out, nil
}
func (u *UseCase) CreateRole(ctx context.Context, actor uuid.UUID, in securityreq.CreateRoleRequestDTO) (*securityresp.RoleResponseDTO, error) {
	var roleOut *entities.Role
	err := u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
		requestedItems, _, err := parseRequestedRolePermissions(ctx, repos, in.Permissions)
		if err != nil {
			return err
		}

		isActive := true
		if in.IsActive != nil {
			isActive = *in.IsActive
		}
		roleType := in.RoleType
		if roleType == "" {
			roleType = enums.RoleTypeSystem
		}
		if roleType != enums.RoleTypeSystem {
			return errUnsupportedRoleType(roleType)
		}

		role := &entities.Role{
			Name:        in.Name,
			Description: in.Description,
			IsActive:    isActive,
			RoleType:    roleType,
			BaseEntity: entities.BaseEntity{
				CreatedAtUtc: time.Now().UTC(),
			},
		}
		if err := repos.Roles().Create(ctx, role); err != nil {
			return err
		}
		rolePerms := cloneRolePermsForRole(role.ID, requestedItems)
		if err := repos.RolePermissions().ReplaceRolePermissions(ctx, role.ID, rolePerms); err != nil {
			return err
		}

		common.CreateAudit(ctx, repos, &actor, "Role", role.ID.String(), "CREATE", role)
		common.CreateAudit(ctx, repos, &actor, "RolePermission", role.ID.String(), "REPLACE", rolePerms)
		roleOut = role
		return nil
	})
	if err != nil {
		return nil, err
	}
	return u.GetRole(ctx, roleOut.ID)
}
func (u *UseCase) UpdateRole(ctx context.Context, actor, roleID uuid.UUID, in securityreq.UpdateRoleRequestDTO) (*securityresp.RoleResponseDTO, error) {
	err := u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
		requestedItems, _, err := parseRequestedRolePermissions(ctx, repos, in.Permissions)
		if err != nil {
			return err
		}

		current, err := repos.Roles().GetByID(ctx, roleID)
		if err != nil {
			return err
		}
		if current.RoleType != enums.RoleTypeSystem {
			return errUnsupportedRoleType(current.RoleType)
		}

		roleType := in.RoleType
		if roleType == "" {
			roleType = current.RoleType
		}
		if roleType != enums.RoleTypeSystem {
			return errUnsupportedRoleType(roleType)
		}

		if in.IsActive != nil {
			current.IsActive = *in.IsActive
		}
		current.Name = in.Name
		current.Description = in.Description
		current.RoleType = roleType
		current.UpdatedAtUtc = time.Now().UTC()
		if err := repos.Roles().Update(ctx, current); err != nil {
			return err
		}

		rolePerms := cloneRolePermsForRole(current.ID, requestedItems)
		if err := repos.RolePermissions().ReplaceRolePermissions(ctx, current.ID, rolePerms); err != nil {
			return err
		}

		common.CreateAudit(ctx, repos, &actor, "Role", current.ID.String(), "UPDATE", current)
		common.CreateAudit(ctx, repos, &actor, "RolePermission", current.ID.String(), "REPLACE", rolePerms)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return u.GetRole(ctx, roleID)
}

func errUnsupportedRoleType(roleType enums.RoleType) error {
	return &unsupportedRoleTypeError{RoleType: roleType}
}

type unsupportedRoleTypeError struct {
	RoleType enums.RoleType
}

func (e *unsupportedRoleTypeError) Error() string {
	return "unsupported role type for create role: " + string(e.RoleType)
}
func (u *UseCase) SetRoleActive(ctx context.Context, actor, id uuid.UUID, active bool) error {
	err := u.uow.Repositories().Roles().SetActive(ctx, id, active)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Role", id.String(), map[bool]string{true: "ACTIVATE", false: "DEACTIVATE"}[active], nil)
	}
	return err
}
func (u *UseCase) RolePermissions(ctx context.Context, roleID uuid.UUID) ([]entities.RoleSubMenuPermission, error) {
	return u.uow.Repositories().RolePermissions().ListRolePermissions(ctx, roleID)
}
func (u *UseCase) ReplaceRolePermissions(ctx context.Context, actor, roleID uuid.UUID, items []entities.RoleSubMenuPermission) error {
	err := u.uow.Repositories().RolePermissions().ReplaceRolePermissions(ctx, roleID, items)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "RolePermission", roleID.String(), "REPLACE", items)
	}
	return err
}

func buildRoleRefMaps(menus []entities.Menu, subMenus []entities.SubMenu, perms []entities.Permission) (map[uuid.UUID]entities.Menu, map[uuid.UUID]entities.SubMenu, map[uuid.UUID]entities.Permission) {
	menuByID := make(map[uuid.UUID]entities.Menu, len(menus))
	for _, m := range menus {
		menuByID[m.ID] = m
	}
	subMenuByID := make(map[uuid.UUID]entities.SubMenu, len(subMenus))
	for _, sm := range subMenus {
		subMenuByID[sm.ID] = sm
	}
	permByID := make(map[uuid.UUID]entities.Permission, len(perms))
	for _, p := range perms {
		permByID[p.ID] = p
	}
	return menuByID, subMenuByID, permByID
}

func buildRoleResponseDTO(
	role entities.Role,
	rolePerms []entities.RoleSubMenuPermission,
	menuByID map[uuid.UUID]entities.Menu,
	subMenuByID map[uuid.UUID]entities.SubMenu,
	permByID map[uuid.UUID]entities.Permission,
) securityresp.RoleResponseDTO {
	menusMap := map[uuid.UUID]*securityresp.RoleMenuDTO{}
	subMenusMap := map[uuid.UUID]*securityresp.RoleSubMenuDTO{}
	for _, rp := range rolePerms {
		sm, okSm := subMenuByID[rp.SubMenuID]
		pm, okPm := permByID[rp.PermissionID]
		if !okSm || !okPm {
			continue
		}
		menu, okMenu := menuByID[sm.MenuID]
		if !okMenu {
			continue
		}

		menuDTO, ok := menusMap[menu.ID]
		if !ok {
			menuDTO = &securityresp.RoleMenuDTO{ID: menu.ID, Code: menu.Code, Name: menu.Name, SubMenus: []securityresp.RoleSubMenuDTO{}}
			menusMap[menu.ID] = menuDTO
		}

		subMenuDTO, ok := subMenusMap[sm.ID]
		if !ok {
			subMenuDTO = &securityresp.RoleSubMenuDTO{ID: sm.ID, Code: sm.Code, Name: sm.Name, Permissions: []securityresp.RolePermissionDTO{}}
			subMenusMap[sm.ID] = subMenuDTO
			menuDTO.SubMenus = append(menuDTO.SubMenus, *subMenuDTO)
		}
		subMenuDTO.Permissions = append(subMenuDTO.Permissions, securityresp.RolePermissionDTO{ID: pm.ID, Code: pm.Code, Name: pm.Name})
	}

	menusList := make([]securityresp.RoleMenuDTO, 0, len(menusMap))
	for _, menu := range menusMap {
		subList := make([]securityresp.RoleSubMenuDTO, 0, len(menu.SubMenus))
		for _, sm := range menu.SubMenus {
			current := subMenusMap[sm.ID]
			if current == nil {
				continue
			}
			sort.Slice(current.Permissions, func(i, j int) bool { return current.Permissions[i].Code < current.Permissions[j].Code })
			subList = append(subList, securityresp.RoleSubMenuDTO{
				ID:          current.ID,
				Code:        current.Code,
				Name:        current.Name,
				Permissions: current.Permissions,
			})
		}
		sort.Slice(subList, func(i, j int) bool { return subList[i].Code < subList[j].Code })
		menusList = append(menusList, securityresp.RoleMenuDTO{
			ID:       menu.ID,
			Code:     menu.Code,
			Name:     menu.Name,
			SubMenus: subList,
		})
	}
	sort.Slice(menusList, func(i, j int) bool { return menusList[i].Code < menusList[j].Code })

	return securityresp.RoleResponseDTO{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsActive:    role.IsActive,
		RoleType:    role.RoleType,
		Menus:       menusList,
	}
}
