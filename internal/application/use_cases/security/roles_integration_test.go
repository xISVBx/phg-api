package security_test

import (
	"context"
	"testing"
	"time"

	securityreq "photogallery/api_go/internal/application/dtos/request/security"
	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/domain/enums"
	"photogallery/api_go/tests/testkit"
)

func TestCreateRole_WithPermissions_Success(t *testing.T) {
	env := testkit.NewIntegrationEnv(t)
	ctx := context.Background()
	now := time.Now().UTC()

	menu := &entities.Menu{
		Code:         "ROLE_CREATE_MENU_" + now.Format("150405"),
		Name:         "Role Create Menu " + now.Format("150405.000000"),
		DisplayOrder: 1,
		IsActive:     true,
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
	}
	if err := env.UoW.Repositories().Menus().Create(ctx, menu); err != nil {
		t.Fatalf("create menu: %v", err)
	}

	subMenu := &entities.SubMenu{
		MenuID:       menu.ID,
		Code:         "ROLE_CREATE_SUBMENU_" + now.Format("150405"),
		Name:         "Role Create SubMenu " + now.Format("150405.000000"),
		DisplayOrder: 1,
		IsActive:     true,
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
	}
	if err := env.UoW.Repositories().SubMenus().Create(ctx, subMenu); err != nil {
		t.Fatalf("create submenu: %v", err)
	}

	permView := &entities.Permission{
		Code: "ROLE_CREATE_PERM_VIEW_" + now.Format("150405"),
		Name: "Role Create Perm View " + now.Format("150405.000000"),
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
	}
	if err := env.UoW.Repositories().Permissions().Create(ctx, permView); err != nil {
		t.Fatalf("create permission view: %v", err)
	}

	permEdit := &entities.Permission{
		Code: "ROLE_CREATE_PERM_EDIT_" + now.Format("150405"),
		Name: "Role Create Perm Edit " + now.Format("150405.000000"),
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
	}
	if err := env.UoW.Repositories().Permissions().Create(ctx, permEdit); err != nil {
		t.Fatalf("create permission edit: %v", err)
	}

	out, err := env.UseCases.Security.CreateRole(ctx, env.ActorUser.ID, securityreq.CreateRoleRequestDTO{
		Name:        "TEST_ROLE_CREATE_" + now.Format("150405.000000"),
		Description: "role created with permissions",
		Permissions: []securityreq.RolePermissionSetItemDTO{
			{
				MenuID:        menu.ID.String(),
				SubMenuID:     subMenu.ID.String(),
				PermissionIDs: []string{permView.ID.String(), permEdit.ID.String()},
			},
		},
	})
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	if out == nil {
		t.Fatalf("expected created role response, got nil")
	}
	if out.RoleType != enums.RoleTypeSystem {
		t.Fatalf("expected role type %s, got %s", enums.RoleTypeSystem, out.RoleType)
	}
	if len(out.Menus) != 1 {
		t.Fatalf("expected one menu in response, got %d", len(out.Menus))
	}
	if len(out.Menus[0].SubMenus) != 1 {
		t.Fatalf("expected one submenu in response, got %d", len(out.Menus[0].SubMenus))
	}
	if len(out.Menus[0].SubMenus[0].Permissions) != 2 {
		t.Fatalf("expected two permissions in response, got %d", len(out.Menus[0].SubMenus[0].Permissions))
	}

	roleID := out.ID
	perms, err := env.UoW.Repositories().RolePermissions().ListRolePermissions(ctx, roleID)
	if err != nil {
		t.Fatalf("list persisted role permissions: %v", err)
	}
	if len(perms) != 2 {
		t.Fatalf("expected 2 persisted role permissions, got %d", len(perms))
	}
}

func TestUpdateRole_WithPermissions_Success(t *testing.T) {
	env := testkit.NewIntegrationEnv(t)
	ctx := context.Background()
	now := time.Now().UTC()

	menu := &entities.Menu{
		Code:         "ROLE_UPDATE_MENU_" + now.Format("150405"),
		Name:         "Role Update Menu " + now.Format("150405.000000"),
		DisplayOrder: 1,
		IsActive:     true,
		BaseEntity:   entities.BaseEntity{CreatedAtUtc: now},
	}
	if err := env.UoW.Repositories().Menus().Create(ctx, menu); err != nil {
		t.Fatalf("create menu: %v", err)
	}

	subMenu := &entities.SubMenu{
		MenuID:       menu.ID,
		Code:         "ROLE_UPDATE_SUBMENU_" + now.Format("150405"),
		Name:         "Role Update SubMenu " + now.Format("150405.000000"),
		DisplayOrder: 1,
		IsActive:     true,
		BaseEntity:   entities.BaseEntity{CreatedAtUtc: now},
	}
	if err := env.UoW.Repositories().SubMenus().Create(ctx, subMenu); err != nil {
		t.Fatalf("create submenu: %v", err)
	}

	permView := &entities.Permission{
		Code:       "ROLE_UPDATE_PERM_VIEW_" + now.Format("150405"),
		Name:       "Role Update Perm View " + now.Format("150405.000000"),
		BaseEntity: entities.BaseEntity{CreatedAtUtc: now},
	}
	if err := env.UoW.Repositories().Permissions().Create(ctx, permView); err != nil {
		t.Fatalf("create permission view: %v", err)
	}
	permEdit := &entities.Permission{
		Code:       "ROLE_UPDATE_PERM_EDIT_" + now.Format("150405"),
		Name:       "Role Update Perm Edit " + now.Format("150405.000000"),
		BaseEntity: entities.BaseEntity{CreatedAtUtc: now},
	}
	if err := env.UoW.Repositories().Permissions().Create(ctx, permEdit); err != nil {
		t.Fatalf("create permission edit: %v", err)
	}

	created, err := env.UseCases.Security.CreateRole(ctx, env.ActorUser.ID, securityreq.CreateRoleRequestDTO{
		Name:        "TEST_ROLE_UPDATE_BASE_" + now.Format("150405.000000"),
		Description: "role before update",
		Permissions: []securityreq.RolePermissionSetItemDTO{
			{
				MenuID:        menu.ID.String(),
				SubMenuID:     subMenu.ID.String(),
				PermissionIDs: []string{permView.ID.String(), permEdit.ID.String()},
			},
		},
	})
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	isActive := false
	updated, err := env.UseCases.Security.UpdateRole(ctx, env.ActorUser.ID, created.ID, securityreq.UpdateRoleRequestDTO{
		Name:        "TEST_ROLE_UPDATED_" + now.Format("150405.000000"),
		Description: "role after update",
		IsActive:    &isActive,
		Permissions: []securityreq.RolePermissionSetItemDTO{
			{
				MenuID:        menu.ID.String(),
				SubMenuID:     subMenu.ID.String(),
				PermissionIDs: []string{permView.ID.String()},
			},
		},
	})
	if err != nil {
		t.Fatalf("update role: %v", err)
	}
	if updated == nil {
		t.Fatalf("expected updated role response, got nil")
	}
	if updated.Name == created.Name {
		t.Fatalf("expected updated role name to change")
	}
	if updated.IsActive {
		t.Fatalf("expected updated role inactive")
	}
	if len(updated.Menus) != 1 || len(updated.Menus[0].SubMenus) != 1 || len(updated.Menus[0].SubMenus[0].Permissions) != 1 {
		t.Fatalf("expected updated role response to include exactly one permission")
	}
	if updated.Menus[0].SubMenus[0].Permissions[0].ID != permView.ID {
		t.Fatalf("expected updated permission %s", permView.ID)
	}

	perms, err := env.UoW.Repositories().RolePermissions().ListRolePermissions(ctx, created.ID)
	if err != nil {
		t.Fatalf("list persisted role permissions: %v", err)
	}
	if len(perms) != 1 {
		t.Fatalf("expected persisted role permissions replaced to 1, got %d", len(perms))
	}
	if perms[0].PermissionID != permView.ID {
		t.Fatalf("expected persisted permission %s, got %s", permView.ID, perms[0].PermissionID)
	}
}
