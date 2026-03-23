package bootstrap

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	appsvc "photogallery/api_go/internal/application/interfaces/services"
	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/domain/enums"
	"photogallery/api_go/internal/infrastructure/config"
)

const advisoryLockKey int64 = 24032301

type Service struct {
	db      *gorm.DB
	cfg     *config.Config
	passSvc appsvc.IPasswordService
}

type Result struct {
	Trigger            string `json:"trigger"`
	Skipped            bool   `json:"skipped"`
	Reason             string `json:"reason,omitempty"`
	CreatedRoles       int    `json:"createdRoles"`
	CreatedPermissions int    `json:"createdPermissions"`
	CreatedMenus       int    `json:"createdMenus"`
	CreatedSubMenus    int    `json:"createdSubMenus"`
	CreatedUsers       int    `json:"createdUsers"`
	CreatedUserRoles   int    `json:"createdUserRoles"`
	CreatedACLGrants   int    `json:"createdACLGrants"`
}

type menuSpec struct {
	Code         string
	Name         string
	DisplayOrder int
}

type subMenuSpec struct {
	MenuCode     string
	Code         string
	Name         string
	Route        string
	DisplayOrder int
}

type permissionSpec struct {
	Code        string
	Name        string
	Description string
}

type roleSpec struct {
	Name        string
	Description string
	RoleType    enums.RoleType
}

type roleGrant struct {
	RoleName        string
	SubMenuCode     string
	PermissionCodes []string
}

var (
	baseRoles = []roleSpec{
		{Name: "Administrador", Description: "Rol administrador base del sistema", RoleType: enums.RoleTypeSystem},
		{Name: "Vendedor", Description: "Rol operativo base para ventas", RoleType: enums.RoleTypeSystem},
	}
	basePermissions = []permissionSpec{
		{Code: "READ", Name: "Read", Description: "Lectura"},
		{Code: "CREATE", Name: "Create", Description: "Creacion"},
		{Code: "UPDATE", Name: "Update", Description: "Actualizacion"},
		{Code: "DELETE", Name: "Delete", Description: "Eliminacion"},
		{Code: "EXECUTE", Name: "Execute", Description: "Ejecucion"},
	}
	baseMenus = []menuSpec{
		{Code: "SECURITY", Name: "Seguridad", DisplayOrder: 1},
		{Code: "OPERATIONS", Name: "Operaciones", DisplayOrder: 2},
		{Code: "ADMIN", Name: "Administracion", DisplayOrder: 3},
	}
	baseSubMenus = []subMenuSpec{
		{MenuCode: "SECURITY", Code: "SEC_USERS", Name: "Usuarios", Route: "/users", DisplayOrder: 1},
		{MenuCode: "SECURITY", Code: "SEC_ROLES", Name: "Roles", Route: "/roles", DisplayOrder: 2},
		{MenuCode: "OPERATIONS", Code: "OPS_CATALOG", Name: "Catalogo", Route: "/products", DisplayOrder: 1},
		{MenuCode: "OPERATIONS", Code: "OPS_SALES", Name: "Ventas", Route: "/sales", DisplayOrder: 2},
		{MenuCode: "OPERATIONS", Code: "OPS_CUSTOMERS", Name: "Clientes", Route: "/customers", DisplayOrder: 3},
		{MenuCode: "OPERATIONS", Code: "OPS_WORKERS", Name: "Trabajadores", Route: "/workers", DisplayOrder: 4},
		{MenuCode: "ADMIN", Code: "ADM_SETTINGS", Name: "Configuracion", Route: "/settings", DisplayOrder: 1},
	}
	baseRoleGrants = []roleGrant{
		{RoleName: "Administrador", SubMenuCode: "SEC_USERS", PermissionCodes: []string{"READ", "CREATE", "UPDATE", "DELETE", "EXECUTE"}},
		{RoleName: "Administrador", SubMenuCode: "SEC_ROLES", PermissionCodes: []string{"READ", "CREATE", "UPDATE", "DELETE", "EXECUTE"}},
		{RoleName: "Administrador", SubMenuCode: "OPS_CATALOG", PermissionCodes: []string{"READ", "CREATE", "UPDATE", "DELETE", "EXECUTE"}},
		{RoleName: "Administrador", SubMenuCode: "OPS_SALES", PermissionCodes: []string{"READ", "CREATE", "UPDATE", "DELETE", "EXECUTE"}},
		{RoleName: "Administrador", SubMenuCode: "OPS_CUSTOMERS", PermissionCodes: []string{"READ", "CREATE", "UPDATE", "DELETE", "EXECUTE"}},
		{RoleName: "Administrador", SubMenuCode: "OPS_WORKERS", PermissionCodes: []string{"READ", "CREATE", "UPDATE", "DELETE", "EXECUTE"}},
		{RoleName: "Administrador", SubMenuCode: "ADM_SETTINGS", PermissionCodes: []string{"READ", "CREATE", "UPDATE", "DELETE", "EXECUTE"}},
		{RoleName: "Vendedor", SubMenuCode: "OPS_CATALOG", PermissionCodes: []string{"READ"}},
		{RoleName: "Vendedor", SubMenuCode: "OPS_SALES", PermissionCodes: []string{"READ", "CREATE", "UPDATE"}},
		{RoleName: "Vendedor", SubMenuCode: "OPS_CUSTOMERS", PermissionCodes: []string{"READ", "CREATE", "UPDATE"}},
		{RoleName: "Vendedor", SubMenuCode: "OPS_WORKERS", PermissionCodes: []string{"READ"}},
	}
)

func NewService(db *gorm.DB, cfg *config.Config, passSvc appsvc.IPasswordService) *Service {
	return &Service{db: db, cfg: cfg, passSvc: passSvc}
}

func (s *Service) RunOnStartup(ctx context.Context) (*Result, error) {
	if !s.cfg.BootstrapAutoEnabled {
		return &Result{Trigger: "startup", Skipped: true, Reason: "BOOTSTRAP_AUTO_ENABLED=false"}, nil
	}
	if !isProdEnv(s.cfg.Env) {
		return &Result{Trigger: "startup", Skipped: true, Reason: "APP_ENV no es prod/production"}, nil
	}
	return s.run(ctx, "startup")
}

func (s *Service) RunManual(ctx context.Context) (*Result, error) {
	if !s.cfg.BootstrapManualEnabled {
		return &Result{Trigger: "manual", Skipped: true, Reason: "BOOTSTRAP_MANUAL_ENABLED=false"}, nil
	}
	return s.run(ctx, "manual")
}

func (s *Service) ValidateManualToken(raw string) bool {
	expected := strings.TrimSpace(s.cfg.BootstrapToken)
	provided := strings.TrimSpace(raw)
	if expected == "" || provided == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expected), []byte(provided)) == 1
}

func (s *Service) run(ctx context.Context, trigger string) (*Result, error) {
	if strings.TrimSpace(s.cfg.BootstrapAdminUser) == "" {
		return nil, errors.New("bootstrap: BOOTSTRAP_ADMIN_USERNAME es obligatorio")
	}
	if strings.TrimSpace(s.cfg.BootstrapAdminEmail) == "" {
		return nil, errors.New("bootstrap: BOOTSTRAP_ADMIN_EMAIL es obligatorio")
	}

	log.Printf("bootstrap: iniciando trigger=%s", trigger)
	result := &Result{Trigger: trigger}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := acquireAdvisoryLock(tx); err != nil {
			return err
		}

		now := time.Now().UTC()
		rolesByName := make(map[string]entities.Role, len(baseRoles))
		for _, spec := range baseRoles {
			role := entities.Role{
				Name:        spec.Name,
				Description: spec.Description,
				IsActive:    true,
				RoleType:    spec.RoleType,
			}
			created, err := ensureRole(tx, &role, now)
			if err != nil {
				return err
			}
			if created {
				result.CreatedRoles++
				log.Printf("bootstrap: rol creado: %s", role.Name)
			} else {
				log.Printf("bootstrap: rol existente: %s", role.Name)
			}
			rolesByName[role.Name] = role
		}

		permissionsByCode := make(map[string]entities.Permission, len(basePermissions))
		for _, spec := range basePermissions {
			perm := entities.Permission{
				Code:        spec.Code,
				Name:        spec.Name,
				Description: spec.Description,
			}
			created, err := ensurePermission(tx, &perm, now)
			if err != nil {
				return err
			}
			if created {
				result.CreatedPermissions++
				log.Printf("bootstrap: permiso creado: %s", perm.Code)
			}
			permissionsByCode[perm.Code] = perm
		}

		menusByCode := make(map[string]entities.Menu, len(baseMenus))
		for _, spec := range baseMenus {
			menu := entities.Menu{
				Code:         spec.Code,
				Name:         spec.Name,
				DisplayOrder: spec.DisplayOrder,
				IsActive:     true,
			}
			created, err := ensureMenu(tx, &menu, now)
			if err != nil {
				return err
			}
			if created {
				result.CreatedMenus++
				log.Printf("bootstrap: menu creado: %s", menu.Code)
			}
			menusByCode[menu.Code] = menu
		}

		subMenusByCode := make(map[string]entities.SubMenu, len(baseSubMenus))
		for _, spec := range baseSubMenus {
			menu, ok := menusByCode[spec.MenuCode]
			if !ok {
				return fmt.Errorf("bootstrap: menu base no encontrado: %s", spec.MenuCode)
			}
			subMenu := entities.SubMenu{
				MenuID:       menu.ID,
				Code:         spec.Code,
				Name:         spec.Name,
				Route:        spec.Route,
				DisplayOrder: spec.DisplayOrder,
				IsActive:     true,
			}
			created, err := ensureSubMenu(tx, &subMenu, now)
			if err != nil {
				return err
			}
			if created {
				result.CreatedSubMenus++
				log.Printf("bootstrap: submenu creado: %s", subMenu.Code)
			}
			subMenusByCode[subMenu.Code] = subMenu
		}

		for _, grant := range baseRoleGrants {
			role, ok := rolesByName[grant.RoleName]
			if !ok {
				return fmt.Errorf("bootstrap: rol de grant no encontrado: %s", grant.RoleName)
			}
			subMenu, ok := subMenusByCode[grant.SubMenuCode]
			if !ok {
				return fmt.Errorf("bootstrap: submenu de grant no encontrado: %s", grant.SubMenuCode)
			}
			for _, permissionCode := range grant.PermissionCodes {
				permission, ok := permissionsByCode[permissionCode]
				if !ok {
					return fmt.Errorf("bootstrap: permiso de grant no encontrado: %s", permissionCode)
				}
				created, err := ensureRoleSubMenuPermission(tx, role.ID, subMenu.ID, permission.ID)
				if err != nil {
					return err
				}
				if created {
					result.CreatedACLGrants++
				}
			}
		}

		adminRole := rolesByName["Administrador"]
		adminUser, userCreated, err := s.ensureAdminUser(tx, now)
		if err != nil {
			return err
		}
		if userCreated {
			result.CreatedUsers++
			log.Printf("bootstrap: usuario admin creado: username=%s email=%s", adminUser.Username, adminUser.Email)
		} else {
			log.Printf("bootstrap: usuario admin existente: username=%s email=%s", adminUser.Username, adminUser.Email)
		}

		userRoleCreated, err := ensurePrimaryUserRole(tx, adminUser.ID, adminRole.ID)
		if err != nil {
			return err
		}
		if userRoleCreated {
			result.CreatedUserRoles++
			log.Printf("bootstrap: rol administrador asignado al usuario %s", adminUser.Username)
		} else {
			log.Printf("bootstrap: el usuario %s ya tiene el rol administrador", adminUser.Username)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	log.Printf(
		"bootstrap: completado trigger=%s roles=%d permisos=%d menus=%d submenus=%d usuarios=%d user_roles=%d grants=%d",
		result.Trigger,
		result.CreatedRoles,
		result.CreatedPermissions,
		result.CreatedMenus,
		result.CreatedSubMenus,
		result.CreatedUsers,
		result.CreatedUserRoles,
		result.CreatedACLGrants,
	)
	return result, nil
}

func (s *Service) ensureAdminUser(tx *gorm.DB, now time.Time) (entities.User, bool, error) {
	var existing entities.User
	err := tx.Where("username = ? OR email = ?", s.cfg.BootstrapAdminUser, s.cfg.BootstrapAdminEmail).First(&existing).Error
	if err == nil {
		if !existing.IsActive {
			if err := tx.Model(&entities.User{}).
				Where("id = ?", existing.ID).
				Updates(map[string]any{
					"is_active":      true,
					"updated_at_utc": now,
				}).Error; err != nil {
				return entities.User{}, false, fmt.Errorf("bootstrap: activar usuario admin existente: %w", err)
			}
			existing.IsActive = true
		}
		return existing, false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.User{}, false, fmt.Errorf("bootstrap: buscar usuario admin: %w", err)
	}

	if strings.TrimSpace(s.cfg.BootstrapAdminPass) == "" {
		return entities.User{}, false, errors.New("bootstrap: BOOTSTRAP_ADMIN_PASSWORD es obligatorio para crear el usuario admin")
	}
	if err := s.passSvc.Validate(s.cfg.BootstrapAdminPass); err != nil {
		return entities.User{}, false, fmt.Errorf("bootstrap: BOOTSTRAP_ADMIN_PASSWORD invalido: %w", err)
	}
	hash, err := s.passSvc.Hash(s.cfg.BootstrapAdminPass)
	if err != nil {
		return entities.User{}, false, fmt.Errorf("bootstrap: hash password admin: %w", err)
	}

	admin := entities.User{
		Username:     s.cfg.BootstrapAdminUser,
		PasswordHash: hash,
		FullName:     s.cfg.BootstrapAdminName,
		Email:        s.cfg.BootstrapAdminEmail,
		IsActive:     true,
	}
	admin.CreatedAtUtc = now
	if err := tx.Create(&admin).Error; err != nil {
		return entities.User{}, false, fmt.Errorf("bootstrap: crear usuario admin: %w", err)
	}
	return admin, true, nil
}

func isProdEnv(env string) bool {
	e := strings.TrimSpace(strings.ToLower(env))
	return e == "prod" || e == "production"
}

func acquireAdvisoryLock(tx *gorm.DB) error {
	var locked bool
	if err := tx.Raw("SELECT pg_try_advisory_xact_lock(?)", advisoryLockKey).Scan(&locked).Error; err != nil {
		return fmt.Errorf("bootstrap: advisory lock: %w", err)
	}
	if !locked {
		return errors.New("bootstrap: otro proceso de bootstrap ya esta en ejecucion")
	}
	return nil
}

func ensureRole(tx *gorm.DB, role *entities.Role, now time.Time) (bool, error) {
	var existing entities.Role
	err := tx.Where("name = ?", role.Name).First(&existing).Error
	if err == nil {
		*role = existing
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("bootstrap: get role %s: %w", role.Name, err)
	}
	role.CreatedAtUtc = now
	if err := tx.Create(role).Error; err != nil {
		return false, fmt.Errorf("bootstrap: create role %s: %w", role.Name, err)
	}
	return true, nil
}

func ensurePermission(tx *gorm.DB, perm *entities.Permission, now time.Time) (bool, error) {
	var existing entities.Permission
	err := tx.Where("code = ?", perm.Code).First(&existing).Error
	if err == nil {
		*perm = existing
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("bootstrap: get permission %s: %w", perm.Code, err)
	}
	perm.CreatedAtUtc = now
	if err := tx.Create(perm).Error; err != nil {
		return false, fmt.Errorf("bootstrap: create permission %s: %w", perm.Code, err)
	}
	return true, nil
}

func ensureMenu(tx *gorm.DB, menu *entities.Menu, now time.Time) (bool, error) {
	var existing entities.Menu
	err := tx.Where("code = ?", menu.Code).First(&existing).Error
	if err == nil {
		*menu = existing
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("bootstrap: get menu %s: %w", menu.Code, err)
	}
	menu.CreatedAtUtc = now
	if err := tx.Create(menu).Error; err != nil {
		return false, fmt.Errorf("bootstrap: create menu %s: %w", menu.Code, err)
	}
	return true, nil
}

func ensureSubMenu(tx *gorm.DB, subMenu *entities.SubMenu, now time.Time) (bool, error) {
	var existing entities.SubMenu
	err := tx.Where("code = ?", subMenu.Code).First(&existing).Error
	if err == nil {
		*subMenu = existing
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("bootstrap: get submenu %s: %w", subMenu.Code, err)
	}
	subMenu.CreatedAtUtc = now
	if err := tx.Create(subMenu).Error; err != nil {
		return false, fmt.Errorf("bootstrap: create submenu %s: %w", subMenu.Code, err)
	}
	return true, nil
}

func ensurePrimaryUserRole(tx *gorm.DB, userID, roleID uuid.UUID) (bool, error) {
	var existing entities.UserRole
	err := tx.Where("user_id = ? AND role_id = ?", userID, roleID).First(&existing).Error
	if err == nil {
		if !existing.IsPrimary {
			if err := tx.Model(&entities.UserRole{}).
				Where("user_id = ? AND role_id = ?", userID, roleID).
				Update("is_primary", true).Error; err != nil {
				return false, fmt.Errorf("bootstrap: update user_role primary: %w", err)
			}
		}
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("bootstrap: get user_role: %w", err)
	}
	if err := tx.Model(&entities.UserRole{}).
		Where("user_id = ? AND is_primary = ?", userID, true).
		Update("is_primary", false).Error; err != nil {
		return false, fmt.Errorf("bootstrap: clear previous primary role: %w", err)
	}
	if err := tx.Create(&entities.UserRole{UserID: userID, RoleID: roleID, IsPrimary: true}).Error; err != nil {
		return false, fmt.Errorf("bootstrap: create user_role: %w", err)
	}
	return true, nil
}

func ensureRoleSubMenuPermission(tx *gorm.DB, roleID, subMenuID, permissionID uuid.UUID) (bool, error) {
	var existing entities.RoleSubMenuPermission
	err := tx.Where("role_id = ? AND sub_menu_id = ? AND permission_id = ?", roleID, subMenuID, permissionID).First(&existing).Error
	if err == nil {
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("bootstrap: get role-submenu-permission: %w", err)
	}
	if err := tx.Create(&entities.RoleSubMenuPermission{
		RoleID:       roleID,
		SubMenuID:    subMenuID,
		PermissionID: permissionID,
	}).Error; err != nil {
		return false, fmt.Errorf("bootstrap: create role-submenu-permission: %w", err)
	}
	return true, nil
}
