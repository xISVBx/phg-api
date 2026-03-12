package security_test

import (
	"context"
	"errors"
	securityreq "photogallery/api_go/internal/application/dtos/request/security"
	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/domain/enums"
	"photogallery/api_go/tests/testkit"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestCreateUser_WithCustomPermissions_Success(t *testing.T) {
	env := testkit.NewIntegrationEnv(t)
	ctx := context.Background()

	now := time.Now().UTC()

	baseRole := &entities.Role{
		Name:        "TEST_ROLE_WITH_CUSTOM_" + now.Format("150405.000000"),
		Description: "role for create user with custom permissions",
		IsActive:    true,
		RoleType:    enums.RoleTypeSystem,
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
	}
	if err := env.UoW.Repositories().Roles().Create(ctx, baseRole); err != nil {
		t.Fatalf("create base role: %v", err)
	}

	menu := &entities.Menu{
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
		Name:     "TEST_MENU_WITH_CUSTOM_" + now.Format("150405.000000"),
		IsActive: true,
	}
	if err := env.UoW.Repositories().Menus().Create(ctx, menu); err != nil {
		t.Fatalf("create menu: %v", err)
	}

	subMenu := &entities.SubMenu{
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
		MenuID:   menu.ID,
		Name:     "TEST_SUBMENU_WITH_CUSTOM_" + now.Format("150405.000000"),
		IsActive: true,
	}
	if err := env.UoW.Repositories().SubMenus().Create(ctx, subMenu); err != nil {
		t.Fatalf("create submenu: %v", err)
	}

	permView := &entities.Permission{
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
		Name: "TEST_PERMISSION_VIEW_WITH_CUSTOM_" + now.Format("150405.000000"),
		Code: "TEST_PERMISSION_VIEW_WITH_CUSTOM_" + now.Format("150405"),
	}
	if err := env.UoW.Repositories().Permissions().Create(ctx, permView); err != nil {
		t.Fatalf("create permission view: %v", err)
	}

	permEdit := &entities.Permission{
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
		Name: "TEST_PERMISSION_EDIT_WITH_CUSTOM_" + now.Format("150405.000000"),
		Code: "TEST_PERMISSION_EDIT_WITH_CUSTOM_" + now.Format("150405"),
	}
	if err := env.UoW.Repositories().Permissions().Create(ctx, permEdit); err != nil {
		t.Fatalf("create permission edit: %v", err)
	}

	baseRolePerms := []entities.RoleSubMenuPermission{
		{
			RoleID:       baseRole.ID,
			SubMenuID:    subMenu.ID,
			PermissionID: permView.ID,
		},
		{
			RoleID:       baseRole.ID,
			SubMenuID:    subMenu.ID,
			PermissionID: permEdit.ID,
		},
	}
	if err := env.UoW.Repositories().RolePermissions().ReplaceRolePermissions(ctx, baseRole.ID, baseRolePerms); err != nil {
		t.Fatalf("replace base role permissions: %v", err)
	}

	var rolesBefore int64
	if err := env.DB.WithContext(ctx).
		Model(&entities.Role{}).
		Count(&rolesBefore).Error; err != nil {
		t.Fatalf("count roles before: %v", err)
	}

	username := "user_with_custom_" + now.Format("150405")

	out, err := env.UseCases.Security.CreateUser(ctx, env.ActorUser.ID, securityreq.CreateUserRequestDTO{
		Username: username,
		Password: "StrongPass123*",
		FullName: "User With Custom Permissions",
		Phone:    "3000000000",
		Email:    "user.with.custom@example.com",
		RoleID:   baseRole.ID.String(),
		Permissions: []securityreq.RolePermissionSetItemDTO{
			{
				MenuID:    menu.ID.String(),
				SubMenuID: subMenu.ID.String(),
				PermissionIDs: []string{
					permView.ID.String(), // solo uno, distinto del rol base
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if out == nil {
		t.Fatalf("expected created user, got nil")
	}

	var persistedUser entities.User
	if err := env.DB.WithContext(ctx).
		Where("id = ?", out.ID).
		First(&persistedUser).Error; err != nil {
		t.Fatalf("find persisted user: %v", err)
	}

	if persistedUser.Username != username {
		t.Fatalf("expected username %q, got %q", username, persistedUser.Username)
	}

	var assignments []entities.UserRole
	if err := env.DB.WithContext(ctx).
		Where("user_id = ?", out.ID).
		Find(&assignments).Error; err != nil {
		t.Fatalf("find user role assignments: %v", err)
	}

	if len(assignments) != 2 {
		t.Fatalf("expected exactly 2 role assignments, got %d", len(assignments))
	}

	var baseAssignment *entities.UserRole
	var customAssignment *entities.UserRole

	for i := range assignments {
		a := assignments[i]
		if a.RoleID == baseRole.ID {
			baseAssignment = &a
		} else {
			customAssignment = &a
		}
	}

	if baseAssignment == nil {
		t.Fatalf("expected assignment for base role %s", baseRole.ID)
	}
	if customAssignment == nil {
		t.Fatalf("expected assignment for custom role")
	}

	if baseAssignment.IsPrimary {
		t.Fatalf("expected base role assignment to be non-primary when custom role exists")
	}

	if !customAssignment.IsPrimary {
		t.Fatalf("expected custom role assignment to be primary")
	}
	primaryCount := 0
	for i := range assignments {
		if assignments[i].IsPrimary {
			primaryCount++
		}
	}
	if primaryCount != 1 {
		t.Fatalf("expected exactly 1 primary role assignment, got %d", primaryCount)
	}

	primaryPerms, err := env.UoW.Repositories().RolePermissions().ListRolePermissions(ctx, customAssignment.RoleID)
	if err != nil {
		t.Fatalf("list primary role permissions: %v", err)
	}
	expectedPrimary := []entities.RoleSubMenuPermission{
		{SubMenuID: subMenu.ID, PermissionID: permView.ID},
	}
	assertRolePermsEqual(t, primaryPerms, expectedPrimary)
	basePerms, err := env.UoW.Repositories().RolePermissions().ListRolePermissions(ctx, baseRole.ID)
	if err != nil {
		t.Fatalf("list base role permissions: %v", err)
	}
	if rolePermSetEqual(primaryPerms, basePerms) {
		t.Fatalf("expected primary(custom) role permissions to differ from base role permissions")
	}

	var rolesAfter int64
	if err := env.DB.WithContext(ctx).
		Model(&entities.Role{}).
		Count(&rolesAfter).Error; err != nil {
		t.Fatalf("count roles after: %v", err)
	}

	if rolesAfter != rolesBefore+1 {
		t.Fatalf("expected one extra custom role created, roles before=%d, roles after=%d", rolesBefore, rolesAfter)
	}

	var customRole entities.Role
	if err := env.DB.WithContext(ctx).
		Where("id = ?", customAssignment.RoleID).
		First(&customRole).Error; err != nil {
		t.Fatalf("find custom role: %v", err)
	}

	if !customRole.IsActive {
		t.Fatalf("expected custom role to be active")
	}
	if customRole.RoleType != enums.RoleTypeCustomUser {
		t.Fatalf("expected custom role type %s, got %s", enums.RoleTypeCustomUser, customRole.RoleType)
	}
	if baseRole.RoleType != enums.RoleTypeSystem {
		t.Fatalf("expected base role type %s, got %s", enums.RoleTypeSystem, baseRole.RoleType)
	}

	var customRolePerms []entities.RoleSubMenuPermission
	if err := env.DB.WithContext(ctx).
		Where("role_id = ?", customRole.ID).
		Find(&customRolePerms).Error; err != nil {
		t.Fatalf("find custom role permissions: %v", err)
	}

	if len(customRolePerms) != 1 {
		t.Fatalf("expected custom role to have exactly 1 permission, got %d", len(customRolePerms))
	}

	if customRolePerms[0].SubMenuID != subMenu.ID {
		t.Fatalf("expected custom role submenu %s, got %s", subMenu.ID, customRolePerms[0].SubMenuID)
	}

	if customRolePerms[0].PermissionID != permView.ID {
		t.Fatalf("expected custom role permission %s, got %s", permView.ID, customRolePerms[0].PermissionID)
	}

	var persistedBaseRolePerms []entities.RoleSubMenuPermission
	if err := env.DB.WithContext(ctx).
		Where("role_id = ?", baseRole.ID).
		Find(&persistedBaseRolePerms).Error; err != nil {
		t.Fatalf("find base role permissions: %v", err)
	}

	if len(persistedBaseRolePerms) != 2 {
		t.Fatalf("expected base role permissions to remain unchanged with 2 items, got %d", len(persistedBaseRolePerms))
	}

}

func TestCreateUser_WithoutCustomPermissions_Success(t *testing.T) {
	env := testkit.NewIntegrationEnv(t)
	ctx := context.Background()

	now := time.Now().UTC()

	baseRole := &entities.Role{
		Name:        "TEST_ROLE_NO_CUSTOM_" + now.Format("150405.000000"),
		Description: "role for create user without custom permissions",
		IsActive:    true,
		RoleType:    enums.RoleTypeSystem,
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
	}

	if err := env.UoW.Repositories().Roles().Create(ctx, baseRole); err != nil {
		t.Fatalf("create base role: %v", err)
	}

	menu := &entities.Menu{
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
		Name:     "TEST_MENU_NO_CUSTOM_" + now.Format("150405.000000"),
		IsActive: true,
	}
	if err := env.UoW.Repositories().Menus().Create(ctx, menu); err != nil {
		t.Fatalf("create menu: %v", err)
	}

	subMenu := &entities.SubMenu{
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
		MenuID:   menu.ID,
		Name:     "TEST_SUBMENU_NO_CUSTOM_" + now.Format("150405.000000"),
		IsActive: true,
	}
	if err := env.UoW.Repositories().SubMenus().Create(ctx, subMenu); err != nil {
		t.Fatalf("create submenu: %v", err)
	}

	permView := &entities.Permission{
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
		Name: "TEST_PERMISSION_VIEW_NO_CUSTOM_" + now.Format("150405.000000"),
		Code: "TEST_PERMISSION_VIEW_NO_CUSTOM_" + now.Format("150405"),
	}
	if err := env.UoW.Repositories().Permissions().Create(ctx, permView); err != nil {
		t.Fatalf("create permission view: %v", err)
	}

	permEdit := &entities.Permission{
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: now,
		},
		Name: "TEST_PERMISSION_EDIT_NO_CUSTOM_" + now.Format("150405.000000"),
		Code: "TEST_PERMISSION_EDIT_NO_CUSTOM_" + now.Format("150405"),
	}
	if err := env.UoW.Repositories().Permissions().Create(ctx, permEdit); err != nil {
		t.Fatalf("create permission edit: %v", err)
	}

	rolePerms := []entities.RoleSubMenuPermission{
		{
			RoleID:       baseRole.ID,
			SubMenuID:    subMenu.ID,
			PermissionID: permView.ID,
		},
		{
			RoleID:       baseRole.ID,
			SubMenuID:    subMenu.ID,
			PermissionID: permEdit.ID,
		},
	}

	if err := env.UoW.Repositories().RolePermissions().ReplaceRolePermissions(ctx, baseRole.ID, rolePerms); err != nil {
		t.Fatalf("replace role submenu permissions: %v", err)
	}

	var rolesBefore int64
	if err := env.DB.WithContext(ctx).
		Model(&entities.Role{}).
		Count(&rolesBefore).Error; err != nil {
		t.Fatalf("count roles before: %v", err)
	}

	username := "user_no_custom_" + now.Format("150405")

	out, err := env.UseCases.Security.CreateUser(ctx, env.ActorUser.ID, securityreq.CreateUserRequestDTO{
		Username: username,
		Password: "StrongPass123*",
		FullName: "User No Custom Permissions",
		Phone:    "3000000000",
		Email:    "user.no.custom@example.com",
		RoleID:   baseRole.ID.String(),
		Permissions: []securityreq.RolePermissionSetItemDTO{
			{
				MenuID:    menu.ID.String(),
				SubMenuID: subMenu.ID.String(),
				PermissionIDs: []string{
					permView.ID.String(),
					permEdit.ID.String(),
				},
			},
		},
	})

	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if out == nil {
		t.Fatalf("expected created user, got nil")
	}

	var persistedUser entities.User
	if err := env.DB.WithContext(ctx).
		Where("id = ?", out.ID).
		First(&persistedUser).Error; err != nil {
		t.Fatalf("find persisted user: %v", err)
	}

	if persistedUser.Username != username {
		t.Fatalf("expected username %q, got %q", username, persistedUser.Username)
	}

	var assignments []entities.UserRole
	if err := env.DB.WithContext(ctx).
		Where("user_id = ?", out.ID).
		Find(&assignments).Error; err != nil {
		t.Fatalf("find user role assignments: %v", err)
	}

	if len(assignments) != 1 {
		t.Fatalf("expected exactly 1 role assignment, got %d", len(assignments))
	}

	if assignments[0].RoleID != baseRole.ID {
		t.Fatalf("expected assigned role %s, got %s", baseRole.ID, assignments[0].RoleID)
	}

	if !assignments[0].IsPrimary {
		t.Fatalf("expected base role assignment to be primary")
	}
	primaryCount := 0
	for i := range assignments {
		if assignments[i].IsPrimary {
			primaryCount++
		}
	}
	if primaryCount != 1 {
		t.Fatalf("expected exactly 1 primary role assignment, got %d", primaryCount)
	}

	primaryPerms, err := env.UoW.Repositories().RolePermissions().ListRolePermissions(ctx, assignments[0].RoleID)
	if err != nil {
		t.Fatalf("list primary role permissions: %v", err)
	}
	assertRolePermsEqual(t, primaryPerms, rolePerms)

	var rolesAfter int64
	if err := env.DB.WithContext(ctx).
		Model(&entities.Role{}).
		Count(&rolesAfter).Error; err != nil {
		t.Fatalf("count roles after: %v", err)
	}

	if rolesAfter != rolesBefore {
		t.Fatalf("expected no extra custom role created, roles before=%d, roles after=%d", rolesBefore, rolesAfter)
	}
}

func TestUpdateUser_WithCustomPermissions_Success(t *testing.T) {
	env := testkit.NewIntegrationEnv(t)
	ctx := context.Background()
	now := time.Now().UTC()

	baseRole := &entities.Role{
		Name:        "TEST_ROLE_UPDATE_WITH_CUSTOM_" + now.Format("150405.000000"),
		Description: "role for update user with custom permissions",
		IsActive:    true,
		RoleType:    enums.RoleTypeSystem,
		BaseEntity:  entities.BaseEntity{CreatedAtUtc: now},
	}
	if err := env.UoW.Repositories().Roles().Create(ctx, baseRole); err != nil {
		t.Fatalf("create base role: %v", err)
	}

	menu := &entities.Menu{
		BaseEntity:   entities.BaseEntity{CreatedAtUtc: now},
		Code:         "TEST_MENU_UPDATE_WITH_CUSTOM_" + now.Format("150405"),
		Name:         "TEST_MENU_UPDATE_WITH_CUSTOM_" + now.Format("150405.000000"),
		DisplayOrder: 1,
		IsActive:     true,
	}
	if err := env.UoW.Repositories().Menus().Create(ctx, menu); err != nil {
		t.Fatalf("create menu: %v", err)
	}

	subMenu := &entities.SubMenu{
		BaseEntity:   entities.BaseEntity{CreatedAtUtc: now},
		MenuID:       menu.ID,
		Code:         "TEST_SUBMENU_UPDATE_WITH_CUSTOM_" + now.Format("150405"),
		Name:         "TEST_SUBMENU_UPDATE_WITH_CUSTOM_" + now.Format("150405.000000"),
		DisplayOrder: 1,
		IsActive:     true,
	}
	if err := env.UoW.Repositories().SubMenus().Create(ctx, subMenu); err != nil {
		t.Fatalf("create submenu: %v", err)
	}

	permView := &entities.Permission{
		BaseEntity: entities.BaseEntity{CreatedAtUtc: now},
		Name:       "TEST_PERMISSION_VIEW_UPDATE_WITH_CUSTOM_" + now.Format("150405.000000"),
		Code:       "TEST_PERMISSION_VIEW_UPDATE_WITH_CUSTOM_" + now.Format("150405"),
	}
	if err := env.UoW.Repositories().Permissions().Create(ctx, permView); err != nil {
		t.Fatalf("create permission view: %v", err)
	}

	permEdit := &entities.Permission{
		BaseEntity: entities.BaseEntity{CreatedAtUtc: now},
		Name:       "TEST_PERMISSION_EDIT_UPDATE_WITH_CUSTOM_" + now.Format("150405.000000"),
		Code:       "TEST_PERMISSION_EDIT_UPDATE_WITH_CUSTOM_" + now.Format("150405"),
	}
	if err := env.UoW.Repositories().Permissions().Create(ctx, permEdit); err != nil {
		t.Fatalf("create permission edit: %v", err)
	}

	baseRolePerms := []entities.RoleSubMenuPermission{
		{RoleID: baseRole.ID, SubMenuID: subMenu.ID, PermissionID: permView.ID},
		{RoleID: baseRole.ID, SubMenuID: subMenu.ID, PermissionID: permEdit.ID},
	}
	if err := env.UoW.Repositories().RolePermissions().ReplaceRolePermissions(ctx, baseRole.ID, baseRolePerms); err != nil {
		t.Fatalf("replace base role permissions: %v", err)
	}

	user := &entities.User{
		Username:     "update_custom_user_" + now.Format("150405"),
		PasswordHash: "hash",
		FullName:     "Update Custom User",
		Phone:        "3000000000",
		Email:        "update.custom.user@example.com",
		IsActive:     true,
	}
	if err := env.UoW.Repositories().Users().Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := env.UoW.Repositories().UserRoles().ReplaceByUser(ctx, user.ID, []entities.UserRole{
		{UserID: user.ID, RoleID: baseRole.ID, IsPrimary: true},
	}); err != nil {
		t.Fatalf("assign base role: %v", err)
	}

	var rolesBefore int64
	if err := env.DB.WithContext(ctx).Model(&entities.Role{}).Count(&rolesBefore).Error; err != nil {
		t.Fatalf("count roles before: %v", err)
	}

	out, err := env.UseCases.Security.UpdateUser(ctx, env.ActorUser.ID, user.ID, securityreq.UpdateUserRequestDTO{
		FullName: "Updated Custom User",
		Phone:    "3000000001",
		Email:    "updated.custom.user@example.com",
		RoleID:   baseRole.ID.String(),
		Permissions: []securityreq.RolePermissionSetItemDTO{
			{
				MenuID:        menu.ID.String(),
				SubMenuID:     subMenu.ID.String(),
				PermissionIDs: []string{permView.ID.String()},
			},
		},
	})
	if err != nil {
		t.Fatalf("update user: %v", err)
	}
	if out == nil {
		t.Fatalf("expected updated user, got nil")
	}

	var assignments []entities.UserRole
	if err := env.DB.WithContext(ctx).Where("user_id = ?", user.ID).Find(&assignments).Error; err != nil {
		t.Fatalf("find user role assignments: %v", err)
	}
	if len(assignments) != 2 {
		t.Fatalf("expected exactly 2 role assignments, got %d", len(assignments))
	}

	var baseAssignment *entities.UserRole
	var customAssignment *entities.UserRole
	primaryCount := 0
	for i := range assignments {
		a := assignments[i]
		if a.RoleID == baseRole.ID {
			baseAssignment = &a
		} else {
			customAssignment = &a
		}
		if a.IsPrimary {
			primaryCount++
		}
	}
	if baseAssignment == nil {
		t.Fatalf("expected assignment for base role %s", baseRole.ID)
	}
	if customAssignment == nil {
		t.Fatalf("expected assignment for custom role")
	}
	if baseAssignment.IsPrimary {
		t.Fatalf("expected base role assignment to be non-primary when custom role exists")
	}
	if !customAssignment.IsPrimary {
		t.Fatalf("expected custom role assignment to be primary")
	}
	if primaryCount != 1 {
		t.Fatalf("expected exactly 1 primary role assignment, got %d", primaryCount)
	}
	var customRole entities.Role
	if err := env.DB.WithContext(ctx).Where("id = ?", customAssignment.RoleID).First(&customRole).Error; err != nil {
		t.Fatalf("find custom role: %v", err)
	}
	if customRole.RoleType != enums.RoleTypeCustomUser {
		t.Fatalf("expected custom role type %s, got %s", enums.RoleTypeCustomUser, customRole.RoleType)
	}

	primaryPerms, err := env.UoW.Repositories().RolePermissions().ListRolePermissions(ctx, customAssignment.RoleID)
	if err != nil {
		t.Fatalf("list primary role permissions: %v", err)
	}
	expectedPrimary := []entities.RoleSubMenuPermission{
		{SubMenuID: subMenu.ID, PermissionID: permView.ID},
	}
	assertRolePermsEqual(t, primaryPerms, expectedPrimary)
	basePerms, err := env.UoW.Repositories().RolePermissions().ListRolePermissions(ctx, baseRole.ID)
	if err != nil {
		t.Fatalf("list base role permissions: %v", err)
	}
	if rolePermSetEqual(primaryPerms, basePerms) {
		t.Fatalf("expected primary(custom) role permissions to differ from base role permissions")
	}

	var rolesAfter int64
	if err := env.DB.WithContext(ctx).Model(&entities.Role{}).Count(&rolesAfter).Error; err != nil {
		t.Fatalf("count roles after: %v", err)
	}
	if rolesAfter != rolesBefore+1 {
		t.Fatalf("expected one extra custom role created, roles before=%d, roles after=%d", rolesBefore, rolesAfter)
	}
}

func TestUpdateUser_WithoutCustomPermissions_Success(t *testing.T) {
	env := testkit.NewIntegrationEnv(t)
	ctx := context.Background()
	now := time.Now().UTC()

	baseRole := &entities.Role{
		Name:        "TEST_ROLE_UPDATE_NO_CUSTOM_" + now.Format("150405.000000"),
		Description: "role for update user without custom permissions",
		IsActive:    true,
		RoleType:    enums.RoleTypeSystem,
		BaseEntity:  entities.BaseEntity{CreatedAtUtc: now},
	}
	if err := env.UoW.Repositories().Roles().Create(ctx, baseRole); err != nil {
		t.Fatalf("create base role: %v", err)
	}

	menu := &entities.Menu{
		BaseEntity:   entities.BaseEntity{CreatedAtUtc: now},
		Code:         "TEST_MENU_UPDATE_NO_CUSTOM_" + now.Format("150405"),
		Name:         "TEST_MENU_UPDATE_NO_CUSTOM_" + now.Format("150405.000000"),
		DisplayOrder: 1,
		IsActive:     true,
	}
	if err := env.UoW.Repositories().Menus().Create(ctx, menu); err != nil {
		t.Fatalf("create menu: %v", err)
	}

	subMenu := &entities.SubMenu{
		BaseEntity:   entities.BaseEntity{CreatedAtUtc: now},
		MenuID:       menu.ID,
		Code:         "TEST_SUBMENU_UPDATE_NO_CUSTOM_" + now.Format("150405"),
		Name:         "TEST_SUBMENU_UPDATE_NO_CUSTOM_" + now.Format("150405.000000"),
		DisplayOrder: 1,
		IsActive:     true,
	}
	if err := env.UoW.Repositories().SubMenus().Create(ctx, subMenu); err != nil {
		t.Fatalf("create submenu: %v", err)
	}

	permView := &entities.Permission{
		BaseEntity: entities.BaseEntity{CreatedAtUtc: now},
		Name:       "TEST_PERMISSION_VIEW_UPDATE_NO_CUSTOM_" + now.Format("150405.000000"),
		Code:       "TEST_PERMISSION_VIEW_UPDATE_NO_CUSTOM_" + now.Format("150405"),
	}
	if err := env.UoW.Repositories().Permissions().Create(ctx, permView); err != nil {
		t.Fatalf("create permission view: %v", err)
	}

	permEdit := &entities.Permission{
		BaseEntity: entities.BaseEntity{CreatedAtUtc: now},
		Name:       "TEST_PERMISSION_EDIT_UPDATE_NO_CUSTOM_" + now.Format("150405.000000"),
		Code:       "TEST_PERMISSION_EDIT_UPDATE_NO_CUSTOM_" + now.Format("150405"),
	}
	if err := env.UoW.Repositories().Permissions().Create(ctx, permEdit); err != nil {
		t.Fatalf("create permission edit: %v", err)
	}

	baseRolePerms := []entities.RoleSubMenuPermission{
		{RoleID: baseRole.ID, SubMenuID: subMenu.ID, PermissionID: permView.ID},
		{RoleID: baseRole.ID, SubMenuID: subMenu.ID, PermissionID: permEdit.ID},
	}
	if err := env.UoW.Repositories().RolePermissions().ReplaceRolePermissions(ctx, baseRole.ID, baseRolePerms); err != nil {
		t.Fatalf("replace base role permissions: %v", err)
	}

	created, err := env.UseCases.Security.CreateUser(ctx, env.ActorUser.ID, securityreq.CreateUserRequestDTO{
		Username: "update_no_custom_user_" + now.Format("150405"),
		Password: "StrongPass123*",
		FullName: "Update No Custom User",
		Phone:    "3000000000",
		Email:    "update.no.custom.user@example.com",
		RoleID:   baseRole.ID.String(),
		Permissions: []securityreq.RolePermissionSetItemDTO{
			{
				MenuID:        menu.ID.String(),
				SubMenuID:     subMenu.ID.String(),
				PermissionIDs: []string{permView.ID.String()},
			},
		},
	})
	if err != nil {
		t.Fatalf("create user with custom role: %v", err)
	}

	var beforeAssignments []entities.UserRole
	if err := env.DB.WithContext(ctx).Where("user_id = ?", created.ID).Find(&beforeAssignments).Error; err != nil {
		t.Fatalf("find assignments before update: %v", err)
	}
	if len(beforeAssignments) != 2 {
		t.Fatalf("expected 2 assignments before update, got %d", len(beforeAssignments))
	}
	var customRoleID string
	for i := range beforeAssignments {
		if beforeAssignments[i].RoleID != baseRole.ID {
			customRoleID = beforeAssignments[i].RoleID.String()
			break
		}
	}
	if customRoleID == "" {
		t.Fatalf("expected custom role assignment before update")
	}

	var rolesBefore int64
	if err := env.DB.WithContext(ctx).Model(&entities.Role{}).Count(&rolesBefore).Error; err != nil {
		t.Fatalf("count roles before update: %v", err)
	}

	_, err = env.UseCases.Security.UpdateUser(ctx, env.ActorUser.ID, created.ID, securityreq.UpdateUserRequestDTO{
		FullName: "Update No Custom User Final",
		Phone:    "3000000010",
		Email:    "update.no.custom.user.final@example.com",
		RoleID:   baseRole.ID.String(),
		Permissions: []securityreq.RolePermissionSetItemDTO{
			{
				MenuID:        menu.ID.String(),
				SubMenuID:     subMenu.ID.String(),
				PermissionIDs: []string{permView.ID.String(), permEdit.ID.String()},
			},
		},
	})
	if err != nil {
		t.Fatalf("update user: %v", err)
	}

	var assignments []entities.UserRole
	if err := env.DB.WithContext(ctx).Where("user_id = ?", created.ID).Find(&assignments).Error; err != nil {
		t.Fatalf("find assignments after update: %v", err)
	}
	if len(assignments) != 1 {
		t.Fatalf("expected 1 assignment after update, got %d", len(assignments))
	}
	if assignments[0].RoleID != baseRole.ID {
		t.Fatalf("expected base role assignment after update, got %s", assignments[0].RoleID)
	}
	if !assignments[0].IsPrimary {
		t.Fatalf("expected base role assignment to be primary after update")
	}

	primaryPerms, err := env.UoW.Repositories().RolePermissions().ListRolePermissions(ctx, assignments[0].RoleID)
	if err != nil {
		t.Fatalf("list primary role permissions: %v", err)
	}
	assertRolePermsEqual(t, primaryPerms, baseRolePerms)

	var rolesAfter int64
	if err := env.DB.WithContext(ctx).Model(&entities.Role{}).Count(&rolesAfter).Error; err != nil {
		t.Fatalf("count roles after update: %v", err)
	}
	if rolesAfter != rolesBefore-1 {
		t.Fatalf("expected custom role deletion, roles before=%d, roles after=%d", rolesBefore, rolesAfter)
	}

	var deletedCustomRole entities.Role
	err = env.DB.WithContext(ctx).Where("id = ?", customRoleID).First(&deletedCustomRole).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected deleted custom role to not exist, err=%v", err)
	}
}

func assertRolePermsEqual(t *testing.T, got, want []entities.RoleSubMenuPermission) {
	t.Helper()
	if !rolePermSetEqual(got, want) {
		t.Fatalf("unexpected role permissions set: got=%v want=%v", toPermKeySet(got), toPermKeySet(want))
	}
}

func rolePermSetEqual(a, b []entities.RoleSubMenuPermission) bool {
	setA := toPermKeySet(a)
	setB := toPermKeySet(b)
	if len(setA) != len(setB) {
		return false
	}
	for k := range setA {
		if _, ok := setB[k]; !ok {
			return false
		}
	}
	return true
}

func toPermKeySet(items []entities.RoleSubMenuPermission) map[string]struct{} {
	out := make(map[string]struct{}, len(items))
	for _, it := range items {
		out[permKey(it.SubMenuID, it.PermissionID)] = struct{}{}
	}
	return out
}

func permKey(subMenuID, permissionID uuid.UUID) string {
	return subMenuID.String() + ":" + permissionID.String()
}

func TestCreateUser_WithoutRole_Fail(t *testing.T) {
	env := testkit.NewIntegrationEnv(t)
	ctx := context.Background()

	username := "no_role_user_" + time.Now().UTC().Format("150405")

	_, err := env.UseCases.Security.CreateUser(ctx, env.ActorUser.ID, securityreq.CreateUserRequestDTO{
		Username:    username,
		Password:    "StrongPass123*",
		FullName:    "No Role User",
		Phone:       "3000000000",
		Email:       "no.role.user@example.com",
		RoleID:      "",
		Permissions: []securityreq.RolePermissionSetItemDTO{},
	})
	if err == nil {
		t.Fatalf("expected role required error")
	}

	var count int64
	if err := env.DB.WithContext(ctx).
		Model(&entities.User{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		t.Fatalf("count users by username: %v", err)
	}

	if count != 0 {
		t.Fatalf("expected no persisted user, got count=%d", count)
	}
}

func TestCreateUser_WithWeakPassword_Fail(t *testing.T) {
	env := testkit.NewIntegrationEnv(t)
	ctx := context.Background()
	baseRole := &entities.Role{
		Name:        "TEST_ROLE_WEAK_PASS_" + time.Now().UTC().Format("150405.000000"),
		Description: "role for weak password test",
		IsActive:    true,
		RoleType:    enums.RoleTypeSystem,
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: time.Now().UTC(),
		},
	}
	if err := env.UoW.Repositories().Roles().Create(ctx, baseRole); err != nil {
		t.Fatalf("create base role: %v", err)
	}
	username := "weak_pass_user_" + time.Now().UTC().Format("150405")
	_, err := env.UseCases.Security.CreateUser(ctx, env.ActorUser.ID, securityreq.CreateUserRequestDTO{
		Username:    username,
		Password:    "123456",
		FullName:    "Weak Password User",
		Phone:       "3000000000",
		Email:       "weak.password.user@example.com",
		RoleID:      baseRole.ID.String(),
		Permissions: []securityreq.RolePermissionSetItemDTO{},
	})
	if err == nil {
		t.Fatalf("expected weak password error")
	}
	var count int64
	if err := env.DB.WithContext(ctx).Model(&entities.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		t.Fatalf("count users by username: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no persisted user, got count=%d", count)
	}
}

func TestCreateUser_WithInvalidEmail_Fail(t *testing.T) {
	env := testkit.NewIntegrationEnv(t)
	ctx := context.Background()

	baseRole := &entities.Role{
		Name:        "TEST_ROLE_INVALID_EMAIL_" + time.Now().UTC().Format("150405.000000"),
		Description: "role for invalid email test",
		IsActive:    true,
		RoleType:    enums.RoleTypeSystem,
		BaseEntity: entities.BaseEntity{
			CreatedAtUtc: time.Now().UTC(),
		},
	}
	if err := env.UoW.Repositories().Roles().Create(ctx, baseRole); err != nil {
		t.Fatalf("create base role: %v", err)
	}

	username := "invalid_email_user_" + time.Now().UTC().Format("150405")

	_, err := env.UseCases.Security.CreateUser(ctx, env.ActorUser.ID, securityreq.CreateUserRequestDTO{
		Username:    username,
		Password:    "StrongPass123*",
		FullName:    "Invalid Email User",
		Phone:       "3000000000",
		Email:       "invalid-email",
		RoleID:      baseRole.ID.String(),
		Permissions: []securityreq.RolePermissionSetItemDTO{},
	})
	if err == nil {
		t.Fatalf("expected invalid email error")
	}

	var count int64
	if err := env.DB.WithContext(ctx).
		Model(&entities.User{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		t.Fatalf("count users by username: %v", err)
	}

	if count != 0 {
		t.Fatalf("expected no persisted user, got count=%d", count)
	}
}
