package security

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	securityreq "photogallery/api_go/internal/application/dtos/request/security"
	securityresp "photogallery/api_go/internal/application/dtos/responses/security"
	appif "photogallery/api_go/internal/application/interfaces"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/domain/enums"
	drepo "photogallery/api_go/internal/domain/repositories"
	valueobjects "photogallery/api_go/internal/domain/value_objects"
)

const customRolePrefix = "CUSTOM_USER_"

func (u *UseCase) ListUsers(ctx context.Context, o drepo.QueryOptions) ([]entities.User, int64, error) {
	return u.uow.Repositories().Users().List(ctx, o)
}

func (u *UseCase) GetUser(ctx context.Context, id uuid.UUID) (*securityresp.UserResponseDTO, error) {
	repos := u.uow.Repositories()
	user, err := repos.Users().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	userRoles, err := repos.UserRoles().ListByUser(ctx, id)
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
	menuByID, subMenuByID, permByID := buildRoleRefMaps(menus, subMenus, perms)

	rolesOut := make([]securityresp.UserRoleDTO, 0, len(userRoles))
	for _, ur := range userRoles {
		role, err := repos.Roles().GetByID(ctx, ur.RoleID)
		if err != nil {
			return nil, err
		}
		rolePerms, err := repos.RolePermissions().ListRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		detailed := buildRoleResponseDTO(*role, rolePerms, menuByID, subMenuByID, permByID)
		rolesOut = append(rolesOut, securityresp.UserRoleDTO{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			IsActive:    role.IsActive,
			IsPrimary:   ur.IsPrimary,
			RoleType:    role.RoleType,
			Menus:       detailed.Menus,
		})
	}
	sort.Slice(rolesOut, func(i, j int) bool {
		if rolesOut[i].IsPrimary != rolesOut[j].IsPrimary {
			return rolesOut[i].IsPrimary
		}
		return rolesOut[i].Name < rolesOut[j].Name
	})

	return &securityresp.UserResponseDTO{
		ID:           user.ID,
		Username:     user.Username,
		FullName:     user.FullName,
		Phone:        user.Phone,
		Email:        user.Email,
		IsActive:     user.IsActive,
		CreatedAtUtc: user.CreatedAtUtc,
		UpdatedAtUtc: user.UpdatedAtUtc,
		Roles:        rolesOut,
	}, nil
}

func (u *UseCase) CreateUser(ctx context.Context, actor uuid.UUID, in securityreq.CreateUserRequestDTO) (*entities.User, error) {
	roleID, err := uuid.Parse(in.RoleID)
	if err != nil {
		return nil, err
	}

	if err := u.pass.Validate(in.Password); err != nil {
		return nil, err
	}

	hash, err := u.pass.Hash(in.Password)
	if err != nil {
		return nil, err
	}

	emailVO, err := valueobjects.NewEmail(in.Email)
	if err != nil {
		return nil, err
	}

	var out *entities.User
	err = u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
		requestedItems, requestedSet, err := parseRequestedRolePermissions(ctx, repos, in.Permissions)
		if err != nil {
			return err
		}

		if _, err := repos.Roles().GetByID(ctx, roleID); err != nil {
			return err
		}

		basePerms, err := repos.RolePermissions().ListRolePermissions(ctx, roleID)
		if err != nil {
			return err
		}
		isSameAsBase := equalRolePermissionSets(requestedSet, rolePermissionSet(basePerms))

		user := &entities.User{
			Username:     in.Username,
			PasswordHash: hash,
			FullName:     in.FullName,
			Phone:        in.Phone,
			Email:        emailVO.String(),
			IsActive:     true,
		}
		if err := repos.Users().Create(ctx, user); err != nil {
			return err
		}

		assignments := []entities.UserRole{{UserID: user.ID, RoleID: roleID, IsPrimary: true}}
		if !isSameAsBase {
			customRole := &entities.Role{
				Name:        customRoleName(user.ID),
				Description: "Rol personalizado para " + user.Username,
				IsActive:    true,
				RoleType:    enums.RoleTypeCustomUser,
			}
			if err := repos.Roles().Create(ctx, customRole); err != nil {
				return err
			}
			customPerms := cloneRolePermsForRole(customRole.ID, requestedItems)
			if err := repos.RolePermissions().ReplaceRolePermissions(ctx, customRole.ID, customPerms); err != nil {
				return err
			}
			assignments = []entities.UserRole{
				{UserID: user.ID, RoleID: roleID, IsPrimary: false},
				{UserID: user.ID, RoleID: customRole.ID, IsPrimary: true},
			}
			common.CreateAudit(ctx, repos, &actor, "Role", customRole.ID.String(), "CREATE_CUSTOM_ROLE", map[string]any{"userId": user.ID, "baseRoleId": roleID})
		}

		if err := repos.UserRoles().ReplaceByUser(ctx, user.ID, assignments); err != nil {
			return err
		}

		common.CreateAudit(ctx, repos, &actor, "User", user.ID.String(), "CREATE", user)
		out = user
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (u *UseCase) UpdateUser(ctx context.Context, actor, id uuid.UUID, in securityreq.UpdateUserRequestDTO) (*entities.User, error) {
	roleID, err := uuid.Parse(in.RoleID)
	if err != nil {
		return nil, err
	}
	var out *entities.User
	err = u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
		requestedItems, requestedSet, err := parseRequestedRolePermissions(ctx, repos, in.Permissions)
		if err != nil {
			return err
		}

		user, err := repos.Users().GetByID(ctx, id)
		if err != nil {
			return err
		}
		if _, err := repos.Roles().GetByID(ctx, roleID); err != nil {
			return err
		}

		basePerms, err := repos.RolePermissions().ListRolePermissions(ctx, roleID)
		if err != nil {
			return err
		}
		isSameAsBase := equalRolePermissionSets(requestedSet, rolePermissionSet(basePerms))

		user.FullName = in.FullName
		user.Phone = in.Phone
		user.Email = in.Email
		if in.IsActive != nil {
			user.IsActive = *in.IsActive
		}
		user.UpdatedAtUtc = time.Now().UTC()
		if err := repos.Users().Update(ctx, user); err != nil {
			return err
		}

		userRoles, err := repos.UserRoles().ListByUser(ctx, user.ID)
		if err != nil {
			return err
		}
		existingCustomRole, err := findCustomRoleByUserRoles(ctx, repos, userRoles)
		if err != nil {
			return err
		}

		if isSameAsBase {
			if err := repos.UserRoles().ReplaceByUser(ctx, user.ID, []entities.UserRole{{UserID: user.ID, RoleID: roleID, IsPrimary: true}}); err != nil {
				return err
			}
			if existingCustomRole != nil {
				if err := deleteCustomRole(ctx, repos, existingCustomRole.ID); err != nil {
					return err
				}
				common.CreateAudit(ctx, repos, &actor, "Role", existingCustomRole.ID.String(), "DELETE_CUSTOM_ROLE", map[string]any{"userId": user.ID, "reason": "permissions_equal_to_base"})
			}
		} else {
			newCustomRole := &entities.Role{
				Name:        customRoleName(user.ID),
				Description: "Rol personalizado para " + user.Username,
				IsActive:    true,
				RoleType:    enums.RoleTypeCustomUser,
			}
			if existingCustomRole != nil && existingCustomRole.Name == newCustomRole.Name {
				newCustomRole.Name = fmt.Sprintf("%s_%d", newCustomRole.Name, time.Now().UTC().Unix())
			}
			if err := repos.Roles().Create(ctx, newCustomRole); err != nil {
				return err
			}
			customPerms := cloneRolePermsForRole(newCustomRole.ID, requestedItems)
			if err := repos.RolePermissions().ReplaceRolePermissions(ctx, newCustomRole.ID, customPerms); err != nil {
				return err
			}
			if err := repos.UserRoles().ReplaceByUser(ctx, user.ID, []entities.UserRole{
				{UserID: user.ID, RoleID: roleID, IsPrimary: false},
				{UserID: user.ID, RoleID: newCustomRole.ID, IsPrimary: true},
			}); err != nil {
				return err
			}
			if existingCustomRole != nil {
				if err := deleteCustomRole(ctx, repos, existingCustomRole.ID); err != nil {
					return err
				}
				common.CreateAudit(ctx, repos, &actor, "Role", existingCustomRole.ID.String(), "DELETE_CUSTOM_ROLE", map[string]any{"userId": user.ID, "reason": "replaced_by_update"})
			}
			common.CreateAudit(ctx, repos, &actor, "Role", newCustomRole.ID.String(), "CREATE_CUSTOM_ROLE", map[string]any{"userId": user.ID, "baseRoleId": roleID})
		}

		common.CreateAudit(ctx, repos, &actor, "User", user.ID.String(), "UPDATE", user)
		out = user
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (u *UseCase) SetUserActive(ctx context.Context, actor, id uuid.UUID, active bool) error {
	if err := u.uow.Repositories().Users().SetActive(ctx, id, active); err != nil {
		return err
	}
	act := "DEACTIVATE"
	if active {
		act = "ACTIVATE"
	}
	common.CreateAudit(ctx, u.uow.Repositories(), &actor, "User", id.String(), act, nil)
	return nil
}

func (u *UseCase) SetAdminPassword(ctx context.Context, actor, id uuid.UUID, newPwd string) error {
	h, err := u.pass.Hash(newPwd)
	if err != nil {
		return err
	}
	if err := u.uow.Repositories().Users().SetPassword(ctx, id, h); err != nil {
		return err
	}
	common.CreateAudit(ctx, u.uow.Repositories(), &actor, "User", id.String(), "SET_PASSWORD", nil)
	return nil
}

func (u *UseCase) GetUserRoles(ctx context.Context, id uuid.UUID) ([]entities.UserRole, error) {
	return u.uow.Repositories().UserRoles().ListByUser(ctx, id)
}

func (u *UseCase) SetUserRole(ctx context.Context, actor, userID, roleID uuid.UUID) error {
	return u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
		if err := repos.UserRoles().SetPrimaryRole(ctx, userID, roleID); err != nil {
			return err
		}
		common.CreateAudit(ctx, repos, &actor, "UserRole", userID.String(), "SET_PRIMARY_ROLE", map[string]any{"roleId": roleID})
		return nil
	})
}

func parseRequestedRolePermissions(ctx context.Context, repos appif.RepositorySet, items []securityreq.RolePermissionSetItemDTO) ([]entities.RoleSubMenuPermission, map[string]struct{}, error) {
	out := make([]entities.RoleSubMenuPermission, 0)
	set := make(map[string]struct{})
	for _, group := range items {
		menuID, err := uuid.Parse(group.MenuID)
		if err != nil {
			return nil, nil, err
		}
		subMenuID, err := uuid.Parse(group.SubMenuID)
		if err != nil {
			return nil, nil, err
		}
		subMenu, err := repos.SubMenus().GetByID(ctx, subMenuID)
		if err != nil {
			return nil, nil, err
		}
		if subMenu.MenuID != menuID {
			return nil, nil, fmt.Errorf("submenu %s no pertenece al menu %s", subMenuID, menuID)
		}
		for _, pidStr := range group.PermissionIDs {
			pid, err := uuid.Parse(pidStr)
			if err != nil {
				return nil, nil, err
			}
			key := subMenuID.String() + ":" + pid.String()
			if _, exists := set[key]; exists {
				continue
			}
			set[key] = struct{}{}
			out = append(out, entities.RoleSubMenuPermission{
				SubMenuID:    subMenuID,
				PermissionID: pid,
			})
		}
	}
	return out, set, nil
}

func rolePermissionSet(items []entities.RoleSubMenuPermission) map[string]struct{} {
	out := make(map[string]struct{}, len(items))
	for _, it := range items {
		out[it.SubMenuID.String()+":"+it.PermissionID.String()] = struct{}{}
	}
	return out
}

func equalRolePermissionSets(a, b map[string]struct{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if _, ok := b[k]; !ok {
			return false
		}
	}
	return true
}

func customRoleName(userID uuid.UUID) string {
	return customRolePrefix + strings.ReplaceAll(userID.String(), "-", "")
}

func findCustomRoleByUserRoles(ctx context.Context, repos appif.RepositorySet, roles []entities.UserRole) (*entities.Role, error) {
	for _, ur := range roles {
		role, err := repos.Roles().GetByID(ctx, ur.RoleID)
		if err != nil {
			return nil, err
		}
		if role.RoleType == enums.RoleTypeCustomUser || strings.HasPrefix(role.Name, customRolePrefix) {
			return role, nil
		}
	}
	return nil, nil
}

func cloneRolePermsForRole(roleID uuid.UUID, items []entities.RoleSubMenuPermission) []entities.RoleSubMenuPermission {
	out := make([]entities.RoleSubMenuPermission, 0, len(items))
	for _, it := range items {
		out = append(out, entities.RoleSubMenuPermission{
			RoleID:       roleID,
			SubMenuID:    it.SubMenuID,
			PermissionID: it.PermissionID,
		})
	}
	return out
}

func deleteCustomRole(ctx context.Context, repos appif.RepositorySet, roleID uuid.UUID) error {
	if err := repos.RolePermissions().ReplaceRolePermissions(ctx, roleID, []entities.RoleSubMenuPermission{}); err != nil {
		return err
	}
	return repos.Roles().Delete(ctx, roleID)
}
