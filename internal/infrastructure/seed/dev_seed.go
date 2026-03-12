package seed

import (
	"context"
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

const (
	seedAdminEmail  = "admin@gmail.com"
	seedSellerEmail = "vendedor@gmail.com"
)

var (
	seedPermissionCodes  = []string{"READ", "CREATE", "UPDATE", "DELETE", "EXECUTE"}
	seedMenuCodes        = []string{"SECURITY", "OPERATIONS", "ADMIN"}
	seedSubMenuCodes     = []string{"SEC_USERS", "SEC_ROLES", "OPS_CATALOG", "OPS_SALES", "OPS_CUSTOMERS", "OPS_WORKERS", "ADM_SETTINGS"}
	seedRoleNames        = []string{"Administrador", "Vendedor"}
	seedCategoryNames    = []string{"Sesiones", "Impresiones"}
	seedProductNames     = []string{"Sesion Basica", "Album Premium", "Foto Enmarcada 20x30"}
	seedWorkerEmails     = []string{"carlos.rojas@gmail.com", "laura.perez@gmail.com", "andres.mejia@gmail.com"}
	seedCustomerCodes    = []string{"CLI-0001", "CLI-0002", "CLI-0003"}
	seedCashCategoryName = []string{"Ventas", "Gastos Operativos"}
	seedCashRefs         = []string{"SEED-CASH-0001", "SEED-CASH-0002"}
	seedAppointmentNotes = []string{"SEED-APPT-0001", "SEED-APPT-0002"}
	seedOrderMarkers     = []string{"SEED-SALE-0001", "SEED-SALE-0002", "SEED-SALE-0003"}
)

func RunDevSeed(ctx context.Context, db *gorm.DB, cfg *config.Config, pass appsvc.IPasswordService) error {
	if !cfg.SeedEnabled {
		return nil
	}
	if !isDevEnv(cfg.Env) {
		log.Printf("seed: SEED_ENABLED=true pero APP_ENV=%s (solo corre en dev/development), omitido", cfg.Env)
		return nil
	}

	var count int64
	if err := db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", seedAdminEmail).Count(&count).Error; err != nil {
		return fmt.Errorf("seed: check admin user by email: %w", err)
	}
	if count > 0 {
		log.Printf("seed: admin '%s' ya existe, se omite seed", seedAdminEmail)
		return nil
	}

	log.Printf("seed: iniciando carga inicial de desarrollo")

	adminHash, err := pass.Hash(cfg.SeedAdminPass)
	if err != nil {
		return fmt.Errorf("seed: hash admin password: %w", err)
	}
	sellerHash, err := pass.Hash(cfg.SeedSellerPass)
	if err != nil {
		return fmt.Errorf("seed: hash seller password: %w", err)
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		adminRole := entities.Role{Name: "Administrador", Description: "Rol administrador del sistema", IsActive: true, RoleType: enums.RoleTypeSystem}
		if err := ensureRole(tx, &adminRole, now); err != nil {
			return err
		}
		sellerRole := entities.Role{Name: "Vendedor", Description: "Rol de ventas", IsActive: true, RoleType: enums.RoleTypeSystem}
		if err := ensureRole(tx, &sellerRole, now); err != nil {
			return err
		}

		adminUser := entities.User{
			Username:     cfg.SeedAdminUser,
			PasswordHash: adminHash,
			FullName:     "Administrador General",
			Email:        seedAdminEmail,
			IsActive:     true,
		}
		if err := ensureUser(tx, &adminUser, now); err != nil {
			return err
		}

		sellerUser := entities.User{
			Username:     cfg.SeedSellerUser,
			PasswordHash: sellerHash,
			FullName:     "Vendedor Demo",
			Email:        seedSellerEmail,
			IsActive:     true,
		}
		if err := ensureUser(tx, &sellerUser, now); err != nil {
			return err
		}

		if err := ensurePrimaryUserRole(tx, adminUser.ID, adminRole.ID); err != nil {
			return err
		}
		if err := ensurePrimaryUserRole(tx, sellerUser.ID, sellerRole.ID); err != nil {
			return err
		}

		permRead := entities.Permission{Code: "READ", Name: "Read", Description: "Lectura"}
		if err := ensurePermission(tx, &permRead, now); err != nil {
			return err
		}
		permCreate := entities.Permission{Code: "CREATE", Name: "Create", Description: "Creación"}
		if err := ensurePermission(tx, &permCreate, now); err != nil {
			return err
		}
		permUpdate := entities.Permission{Code: "UPDATE", Name: "Update", Description: "Actualización"}
		if err := ensurePermission(tx, &permUpdate, now); err != nil {
			return err
		}
		permDelete := entities.Permission{Code: "DELETE", Name: "Delete", Description: "Eliminación"}
		if err := ensurePermission(tx, &permDelete, now); err != nil {
			return err
		}
		permExecute := entities.Permission{Code: "EXECUTE", Name: "Execute", Description: "Ejecución"}
		if err := ensurePermission(tx, &permExecute, now); err != nil {
			return err
		}

		securityMenu := entities.Menu{Code: "SECURITY", Name: "Seguridad", DisplayOrder: 1, IsActive: true}
		if err := ensureMenu(tx, &securityMenu, now); err != nil {
			return err
		}
		opsMenu := entities.Menu{Code: "OPERATIONS", Name: "Operaciones", DisplayOrder: 2, IsActive: true}
		if err := ensureMenu(tx, &opsMenu, now); err != nil {
			return err
		}
		adminMenu := entities.Menu{Code: "ADMIN", Name: "Administración", DisplayOrder: 3, IsActive: true}
		if err := ensureMenu(tx, &adminMenu, now); err != nil {
			return err
		}

		subUsers := entities.SubMenu{MenuID: securityMenu.ID, Code: "SEC_USERS", Name: "Usuarios", Route: "/users", DisplayOrder: 1, IsActive: true}
		if err := ensureSubMenu(tx, &subUsers, now); err != nil {
			return err
		}
		subRoles := entities.SubMenu{MenuID: securityMenu.ID, Code: "SEC_ROLES", Name: "Roles", Route: "/roles", DisplayOrder: 2, IsActive: true}
		if err := ensureSubMenu(tx, &subRoles, now); err != nil {
			return err
		}
		subCatalog := entities.SubMenu{MenuID: opsMenu.ID, Code: "OPS_CATALOG", Name: "Catálogo", Route: "/products", DisplayOrder: 1, IsActive: true}
		if err := ensureSubMenu(tx, &subCatalog, now); err != nil {
			return err
		}
		subSales := entities.SubMenu{MenuID: opsMenu.ID, Code: "OPS_SALES", Name: "Ventas", Route: "/sales", DisplayOrder: 2, IsActive: true}
		if err := ensureSubMenu(tx, &subSales, now); err != nil {
			return err
		}
		subCustomers := entities.SubMenu{MenuID: opsMenu.ID, Code: "OPS_CUSTOMERS", Name: "Clientes", Route: "/customers", DisplayOrder: 3, IsActive: true}
		if err := ensureSubMenu(tx, &subCustomers, now); err != nil {
			return err
		}
		subWorkers := entities.SubMenu{MenuID: opsMenu.ID, Code: "OPS_WORKERS", Name: "Trabajadores", Route: "/workers", DisplayOrder: 4, IsActive: true}
		if err := ensureSubMenu(tx, &subWorkers, now); err != nil {
			return err
		}
		subSettings := entities.SubMenu{MenuID: adminMenu.ID, Code: "ADM_SETTINGS", Name: "Configuración", Route: "/settings", DisplayOrder: 1, IsActive: true}
		if err := ensureSubMenu(tx, &subSettings, now); err != nil {
			return err
		}

		adminSubMenus := []entities.SubMenu{subUsers, subRoles, subCatalog, subSales, subCustomers, subWorkers, subSettings}
		adminPerms := []entities.Permission{permRead, permCreate, permUpdate, permDelete, permExecute}
		for _, sm := range adminSubMenus {
			for _, pm := range adminPerms {
				if err := ensureRoleSubMenuPermission(tx, adminRole.ID, sm.ID, pm.ID); err != nil {
					return err
				}
			}
		}

		catStudio := entities.Category{Name: "Sesiones", Description: "Servicios de sesiones fotográficas", IsActive: true}
		if err := ensureCategory(tx, &catStudio, now); err != nil {
			return err
		}
		catPrint := entities.Category{Name: "Impresiones", Description: "Productos de impresión", IsActive: true}
		if err := ensureCategory(tx, &catPrint, now); err != nil {
			return err
		}

		lead := 3
		studioProduct := entities.Product{
			CategoryID:       catStudio.ID,
			Name:             "Sesion Basica",
			Type:             "Service",
			BasePrice:        120,
			Cost:             30,
			CommissionType:   "Percentage",
			CommissionValue:  10,
			RequiresDelivery: false,
			DefaultLeadDays:  nil,
			IsActive:         true,
			Notes:            "Seed dev",
		}
		if err := ensureProduct(tx, &studioProduct, now); err != nil {
			return err
		}
		printProduct := entities.Product{
			CategoryID:       catPrint.ID,
			Name:             "Album Premium",
			Type:             "Product",
			BasePrice:        250,
			Cost:             120,
			CommissionType:   "Fixed",
			CommissionValue:  25,
			RequiresDelivery: true,
			DefaultLeadDays:  &lead,
			IsActive:         true,
			Notes:            "Seed dev",
		}
		if err := ensureProduct(tx, &printProduct, now); err != nil {
			return err
		}
		frameProduct := entities.Product{
			CategoryID:       catPrint.ID,
			Name:             "Foto Enmarcada 20x30",
			Type:             "Product",
			BasePrice:        90,
			Cost:             35,
			CommissionType:   "Fixed",
			CommissionValue:  8,
			RequiresDelivery: false,
			IsActive:         true,
			Notes:            "Seed dev",
		}
		if err := ensureProduct(tx, &frameProduct, now); err != nil {
			return err
		}

		workerA := entities.Worker{
			FullName:     "Carlos Rojas",
			Phone:        "+57 300 111 2233",
			Email:        "carlos.rojas@gmail.com",
			IsActive:     true,
			FixedSalary:  1500,
			SalaryPeriod: "Monthly",
			Notes:        "Seed dev",
		}
		if err := ensureWorker(tx, &workerA, now); err != nil {
			return err
		}
		workerB := entities.Worker{
			FullName:     "Laura Perez",
			Phone:        "+57 300 222 3344",
			Email:        "laura.perez@gmail.com",
			IsActive:     true,
			FixedSalary:  1650,
			SalaryPeriod: "Monthly",
			Notes:        "Seed dev",
		}
		if err := ensureWorker(tx, &workerB, now); err != nil {
			return err
		}
		workerC := entities.Worker{
			FullName:     "Andres Mejia",
			Phone:        "+57 300 333 4455",
			Email:        "andres.mejia@gmail.com",
			IsActive:     true,
			FixedSalary:  1400,
			SalaryPeriod: "Monthly",
			Notes:        "Seed dev",
		}
		if err := ensureWorker(tx, &workerC, now); err != nil {
			return err
		}

		customerA := entities.Customer{
			FullName:     "Ana Torres",
			Phone:        "+57 301 000 1111",
			Email:        "ana.torres@gmail.com",
			CustomerCode: "CLI-0001",
			Document:     "CC-1001",
			Notes:        "Seed dev",
			IsActive:     true,
		}
		if err := ensureCustomer(tx, &customerA, now); err != nil {
			return err
		}
		customerB := entities.Customer{
			FullName:     "Miguel Herrera",
			Phone:        "+57 302 000 2222",
			Email:        "miguel.herrera@gmail.com",
			CustomerCode: "CLI-0002",
			Document:     "CC-1002",
			Notes:        "Seed dev",
			IsActive:     true,
		}
		if err := ensureCustomer(tx, &customerB, now); err != nil {
			return err
		}
		customerC := entities.Customer{
			FullName:     "Sofia Ramirez",
			Phone:        "+57 303 000 3333",
			Email:        "sofia.ramirez@gmail.com",
			CustomerCode: "CLI-0003",
			Document:     "CC-1003",
			Notes:        "Seed dev",
			IsActive:     true,
		}
		if err := ensureCustomer(tx, &customerC, now); err != nil {
			return err
		}

		if err := ensureSaleWithOrder(
			tx,
			"SEED-SALE-0001",
			sellerUser.ID,
			&customerA.ID,
			[]entities.Product{studioProduct, printProduct},
			now,
		); err != nil {
			return err
		}
		if err := ensureSaleWithOrder(
			tx,
			"SEED-SALE-0002",
			sellerUser.ID,
			&customerB.ID,
			[]entities.Product{studioProduct, frameProduct},
			now.Add(-48*time.Hour),
		); err != nil {
			return err
		}
		if err := ensureSaleWithOrder(
			tx,
			"SEED-SALE-0003",
			sellerUser.ID,
			&customerC.ID,
			[]entities.Product{frameProduct},
			now.Add(-24*time.Hour),
		); err != nil {
			return err
		}

		catIncome := entities.CashCategory{Name: "Ventas", Type: "Income", IsActive: true}
		if err := ensureCashCategory(tx, &catIncome, now); err != nil {
			return err
		}
		catExpense := entities.CashCategory{Name: "Gastos Operativos", Type: "Expense", IsActive: true}
		if err := ensureCashCategory(tx, &catExpense, now); err != nil {
			return err
		}
		if err := ensureCashMovement(tx, "SEED-CASH-0001", &catIncome, "Income", "Cash", 370, sellerUser.ID, now); err != nil {
			return err
		}
		if err := ensureCashMovement(tx, "SEED-CASH-0002", &catExpense, "Expense", "Transfer", 120, adminUser.ID, now.Add(-12*time.Hour)); err != nil {
			return err
		}

		if err := ensureAppointment(tx, "SEED-APPT-0001", customerA.ID, studioProduct.ID, sellerUser.ID, now.Add(24*time.Hour)); err != nil {
			return err
		}
		if err := ensureAppointment(tx, "SEED-APPT-0002", customerB.ID, studioProduct.ID, sellerUser.ID, now.Add(48*time.Hour)); err != nil {
			return err
		}

		log.Printf("seed: carga inicial completada (usuarios, catalogo, trabajadores, pedidos, citas y caja)")
		return nil
	})
}

func ResetDevSeedData(ctx context.Context, db *gorm.DB, cfg *config.Config) error {
	if !isDevEnv(cfg.Env) {
		return fmt.Errorf("seed reset permitido solo en dev/development (APP_ENV=%s)", cfg.Env)
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var saleItemSaleIDs []uuid.UUID
		if err := tx.Model(&entities.SaleItem{}).
			Where("notes IN ?", seedOrderMarkers).
			Distinct().
			Pluck("sale_id", &saleItemSaleIDs).Error; err != nil {
			return fmt.Errorf("seed reset: read sale ids: %w", err)
		}

		var workOrderIDs []uuid.UUID
		if err := tx.Model(&entities.WorkOrder{}).
			Where("notes IN ?", seedOrderMarkers).
			Pluck("id", &workOrderIDs).Error; err != nil {
			return fmt.Errorf("seed reset: read work order ids: %w", err)
		}

		if len(workOrderIDs) > 0 {
			if err := tx.Where("work_order_id IN ?", workOrderIDs).Delete(&entities.WorkOrderItem{}).Error; err != nil {
				return fmt.Errorf("seed reset: delete work order items: %w", err)
			}
			if err := tx.Where("id IN ?", workOrderIDs).Delete(&entities.WorkOrder{}).Error; err != nil {
				return fmt.Errorf("seed reset: delete work orders: %w", err)
			}
		}

		if len(saleItemSaleIDs) > 0 {
			if err := tx.Where("sale_id IN ?", saleItemSaleIDs).Delete(&entities.SalePayment{}).Error; err != nil {
				return fmt.Errorf("seed reset: delete sale payments: %w", err)
			}
			if err := tx.Where("sale_id IN ?", saleItemSaleIDs).Delete(&entities.SaleItem{}).Error; err != nil {
				return fmt.Errorf("seed reset: delete sale items by sale_id: %w", err)
			}
			if err := tx.Where("id IN ?", saleItemSaleIDs).Delete(&entities.Sale{}).Error; err != nil {
				return fmt.Errorf("seed reset: delete sales: %w", err)
			}
		} else if err := tx.Where("notes IN ?", seedOrderMarkers).Delete(&entities.SaleItem{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete sale items by marker: %w", err)
		}

		if err := tx.Where("notes IN ?", seedAppointmentNotes).Delete(&entities.Appointment{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete appointments: %w", err)
		}

		if err := tx.Where("reference IN ?", seedCashRefs).Delete(&entities.CashMovement{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete cash movements: %w", err)
		}
		if err := tx.Where("name IN ?", seedCashCategoryName).Delete(&entities.CashCategory{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete cash categories: %w", err)
		}

		if err := tx.Where("email IN ?", seedWorkerEmails).Delete(&entities.Worker{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete workers: %w", err)
		}
		if err := tx.Where("customer_code IN ?", seedCustomerCodes).Delete(&entities.Customer{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete customers: %w", err)
		}
		if err := tx.Where("name IN ?", seedProductNames).Delete(&entities.Product{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete products: %w", err)
		}
		if err := tx.Where("name IN ?", seedCategoryNames).Delete(&entities.Category{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete categories: %w", err)
		}

		var seedUserIDs []uuid.UUID
		if err := tx.Model(&entities.User{}).
			Where("email IN ?", []string{seedAdminEmail, seedSellerEmail}).
			Pluck("id", &seedUserIDs).Error; err != nil {
			return fmt.Errorf("seed reset: read user ids: %w", err)
		}
		if len(seedUserIDs) > 0 {
			if err := tx.Where("user_id IN ?", seedUserIDs).Delete(&entities.UserPermissionOverride{}).Error; err != nil {
				return fmt.Errorf("seed reset: delete user overrides: %w", err)
			}
			if err := tx.Where("user_id IN ?", seedUserIDs).Delete(&entities.UserRole{}).Error; err != nil {
				return fmt.Errorf("seed reset: delete user roles: %w", err)
			}
		}

		if err := tx.Where("email IN ?", []string{seedAdminEmail, seedSellerEmail}).Delete(&entities.User{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete users: %w", err)
		}

		var seedRoleIDs []uuid.UUID
		if err := tx.Model(&entities.Role{}).Where("name IN ?", seedRoleNames).Pluck("id", &seedRoleIDs).Error; err != nil {
			return fmt.Errorf("seed reset: read role ids: %w", err)
		}
		if len(seedRoleIDs) > 0 {
			if err := tx.Where("role_id IN ?", seedRoleIDs).Delete(&entities.RoleSubMenuPermission{}).Error; err != nil {
				return fmt.Errorf("seed reset: delete role permissions: %w", err)
			}
		}

		if err := tx.Where("name IN ?", seedRoleNames).Delete(&entities.Role{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete roles: %w", err)
		}
		if err := tx.Where("code IN ?", seedSubMenuCodes).Delete(&entities.SubMenu{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete submenus: %w", err)
		}
		if err := tx.Where("code IN ?", seedMenuCodes).Delete(&entities.Menu{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete menus: %w", err)
		}
		if err := tx.Where("code IN ?", seedPermissionCodes).Delete(&entities.Permission{}).Error; err != nil {
			return fmt.Errorf("seed reset: delete permissions: %w", err)
		}
		return nil
	})
}

func isDevEnv(env string) bool {
	e := strings.TrimSpace(strings.ToLower(env))
	return e == "dev" || e == "development"
}

func ensureRole(tx *gorm.DB, role *entities.Role, now time.Time) error {
	var existing entities.Role
	err := tx.Where("name = ?", role.Name).First(&existing).Error
	if err == nil {
		*role = existing
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get role %s: %w", role.Name, err)
	}
	role.CreatedAtUtc = now
	if err := tx.Create(role).Error; err != nil {
		return fmt.Errorf("seed: create role %s: %w", role.Name, err)
	}
	return nil
}

func ensureUser(tx *gorm.DB, user *entities.User, now time.Time) error {
	var existing entities.User
	err := tx.Where("username = ?", user.Username).First(&existing).Error
	if err == nil {
		*user = existing
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get user %s: %w", user.Username, err)
	}
	user.CreatedAtUtc = now
	if err := tx.Create(user).Error; err != nil {
		return fmt.Errorf("seed: create user %s: %w", user.Username, err)
	}
	return nil
}

func ensurePrimaryUserRole(tx *gorm.DB, userID, roleID uuid.UUID) error {
	var rel entities.UserRole
	err := tx.Where("user_id = ? AND role_id = ?", userID, roleID).First(&rel).Error
	if err == nil {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get user_role: %w", err)
	}
	if err := tx.Create(&entities.UserRole{UserID: userID, RoleID: roleID, IsPrimary: true}).Error; err != nil {
		return fmt.Errorf("seed: create user_role: %w", err)
	}
	return nil
}

func ensureCategory(tx *gorm.DB, cat *entities.Category, now time.Time) error {
	var existing entities.Category
	err := tx.Where("name = ?", cat.Name).First(&existing).Error
	if err == nil {
		*cat = existing
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get category %s: %w", cat.Name, err)
	}
	cat.CreatedAtUtc = now
	if err := tx.Create(cat).Error; err != nil {
		return fmt.Errorf("seed: create category %s: %w", cat.Name, err)
	}
	return nil
}

func ensureMenu(tx *gorm.DB, menu *entities.Menu, now time.Time) error {
	var existing entities.Menu
	err := tx.Where("code = ?", menu.Code).First(&existing).Error
	if err == nil {
		*menu = existing
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get menu %s: %w", menu.Code, err)
	}
	menu.CreatedAtUtc = now
	if err := tx.Create(menu).Error; err != nil {
		return fmt.Errorf("seed: create menu %s: %w", menu.Code, err)
	}
	return nil
}

func ensureSubMenu(tx *gorm.DB, sub *entities.SubMenu, now time.Time) error {
	var existing entities.SubMenu
	err := tx.Where("code = ?", sub.Code).First(&existing).Error
	if err == nil {
		*sub = existing
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get submenu %s: %w", sub.Code, err)
	}
	sub.CreatedAtUtc = now
	if err := tx.Create(sub).Error; err != nil {
		return fmt.Errorf("seed: create submenu %s: %w", sub.Code, err)
	}
	return nil
}

func ensurePermission(tx *gorm.DB, perm *entities.Permission, now time.Time) error {
	var existing entities.Permission
	err := tx.Where("code = ?", perm.Code).First(&existing).Error
	if err == nil {
		*perm = existing
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get permission %s: %w", perm.Code, err)
	}
	perm.CreatedAtUtc = now
	if err := tx.Create(perm).Error; err != nil {
		return fmt.Errorf("seed: create permission %s: %w", perm.Code, err)
	}
	return nil
}

func ensureRoleSubMenuPermission(tx *gorm.DB, roleID, subMenuID, permissionID uuid.UUID) error {
	var existing entities.RoleSubMenuPermission
	err := tx.Where("role_id = ? AND sub_menu_id = ? AND permission_id = ?", roleID, subMenuID, permissionID).First(&existing).Error
	if err == nil {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get role-submenu-permission: %w", err)
	}
	if err := tx.Create(&entities.RoleSubMenuPermission{
		RoleID:       roleID,
		SubMenuID:    subMenuID,
		PermissionID: permissionID,
	}).Error; err != nil {
		return fmt.Errorf("seed: create role-submenu-permission: %w", err)
	}
	return nil
}

func ensureProduct(tx *gorm.DB, prod *entities.Product, now time.Time) error {
	var existing entities.Product
	err := tx.Where("name = ? AND category_id = ?", prod.Name, prod.CategoryID).First(&existing).Error
	if err == nil {
		*prod = existing
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get product %s: %w", prod.Name, err)
	}
	prod.CreatedAtUtc = now
	if err := tx.Create(prod).Error; err != nil {
		return fmt.Errorf("seed: create product %s: %w", prod.Name, err)
	}
	return nil
}

func ensureWorker(tx *gorm.DB, worker *entities.Worker, now time.Time) error {
	var existing entities.Worker
	err := tx.Where("email = ?", worker.Email).First(&existing).Error
	if err == nil {
		*worker = existing
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get worker %s: %w", worker.Email, err)
	}
	worker.CreatedAtUtc = now
	if err := tx.Create(worker).Error; err != nil {
		return fmt.Errorf("seed: create worker %s: %w", worker.Email, err)
	}
	return nil
}

func ensureCustomer(tx *gorm.DB, customer *entities.Customer, now time.Time) error {
	var existing entities.Customer
	err := tx.Where("customer_code = ?", customer.CustomerCode).First(&existing).Error
	if err == nil {
		*customer = existing
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get customer %s: %w", customer.CustomerCode, err)
	}
	customer.CreatedAtUtc = now
	if err := tx.Create(customer).Error; err != nil {
		return fmt.Errorf("seed: create customer %s: %w", customer.CustomerCode, err)
	}
	return nil
}

func ensureSaleWithOrder(
	tx *gorm.DB,
	marker string,
	sellerUserID uuid.UUID,
	customerID *uuid.UUID,
	products []entities.Product,
	now time.Time,
) error {
	var existingWO entities.WorkOrder
	err := tx.Where("notes = ?", marker).First(&existingWO).Error
	if err == nil {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get seeded order %s: %w", marker, err)
	}
	if len(products) == 0 {
		return fmt.Errorf("seed: order %s requires at least one product", marker)
	}

	sale := entities.Sale{
		CustomerID:   customerID,
		SellerUserID: sellerUserID,
		Status:       string(enums.SaleStatusPending),
	}
	items := make([]entities.SaleItem, 0, len(products))
	var subtotal, totalCost, totalCommission float64
	requiresDelivery := false
	for _, p := range products {
		subtotal += p.BasePrice
		totalCost += p.Cost
		if strings.EqualFold(p.CommissionType, string(enums.CommissionPercent)) {
			totalCommission += p.BasePrice * (p.CommissionValue / 100)
		}
		if strings.EqualFold(p.CommissionType, string(enums.CommissionFixed)) {
			totalCommission += p.CommissionValue
		}
		items = append(items, entities.SaleItem{
			ProductID:                p.ID,
			Quantity:                 1,
			UnitPriceSnapshot:        p.BasePrice,
			UnitCostSnapshot:         p.Cost,
			CommissionTypeSnapshot:   p.CommissionType,
			CommissionValueSnapshot:  p.CommissionValue,
			RequiresDeliverySnapshot: p.RequiresDelivery,
			LeadDaysSnapshot:         p.DefaultLeadDays,
			Notes:                    marker,
		})
		requiresDelivery = requiresDelivery || p.RequiresDelivery
	}
	sale.Subtotal = subtotal
	sale.Total = subtotal
	sale.TotalCostSnapshot = totalCost
	sale.TotalCommissionSnapshot = totalCommission
	sale.CreatedAtUtc = now
	if err := tx.Create(&sale).Error; err != nil {
		return fmt.Errorf("seed: create sale %s: %w", marker, err)
	}
	for i := range items {
		items[i].SaleID = sale.ID
		items[i].CreatedAtUtc = now
	}
	if err := tx.Create(&items).Error; err != nil {
		return fmt.Errorf("seed: create sale items %s: %w", marker, err)
	}

	if !requiresDelivery {
		return nil
	}

	due := now.AddDate(0, 0, 3)
	wo := entities.WorkOrder{
		SaleID:     sale.ID,
		Status:     string(enums.WorkOrderCreated),
		DueDateUtc: &due,
		Notes:      marker,
	}
	wo.CreatedAtUtc = now
	if err := tx.Create(&wo).Error; err != nil {
		return fmt.Errorf("seed: create work order %s: %w", marker, err)
	}

	woItems := make([]entities.WorkOrderItem, 0, len(items))
	for _, it := range items {
		woItems = append(woItems, entities.WorkOrderItem{
			WorkOrderID: wo.ID,
			SaleItemID:  it.ID,
			Status:      string(enums.WorkOrderCreated),
			DueDateUtc:  &due,
			Notes:       marker,
			BaseEntity: entities.BaseEntity{
				CreatedAtUtc: now,
			},
		})
	}
	if err := tx.Create(&woItems).Error; err != nil {
		return fmt.Errorf("seed: create work order items %s: %w", marker, err)
	}
	return nil
}

func ensureCashCategory(tx *gorm.DB, cat *entities.CashCategory, now time.Time) error {
	var existing entities.CashCategory
	err := tx.Where("name = ?", cat.Name).First(&existing).Error
	if err == nil {
		*cat = existing
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get cash category %s: %w", cat.Name, err)
	}
	cat.CreatedAtUtc = now
	if err := tx.Create(cat).Error; err != nil {
		return fmt.Errorf("seed: create cash category %s: %w", cat.Name, err)
	}
	return nil
}

func ensureCashMovement(
	tx *gorm.DB,
	marker string,
	category *entities.CashCategory,
	kind string,
	method string,
	amount float64,
	actor uuid.UUID,
	now time.Time,
) error {
	var existing entities.CashMovement
	err := tx.Where("reference = ?", marker).First(&existing).Error
	if err == nil {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get cash movement %s: %w", marker, err)
	}
	m := entities.CashMovement{
		Type:              kind,
		CategoryID:        category.ID,
		Method:            method,
		Amount:            amount,
		Reference:         marker,
		RelatedEntityType: "Seed",
		RelatedEntityID:   marker,
		Notes:             "Seed dev",
		CreatedByUserID:   actor,
	}
	m.CreatedAtUtc = now
	if err := tx.Create(&m).Error; err != nil {
		return fmt.Errorf("seed: create cash movement %s: %w", marker, err)
	}
	return nil
}

func ensureAppointment(
	tx *gorm.DB,
	marker string,
	customerID uuid.UUID,
	productID uuid.UUID,
	actor uuid.UUID,
	startAt time.Time,
) error {
	var existing entities.Appointment
	err := tx.Where("notes = ?", marker).First(&existing).Error
	if err == nil {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("seed: get appointment %s: %w", marker, err)
	}
	end := startAt.Add(2 * time.Hour)
	app := entities.Appointment{
		CustomerID:      customerID,
		ProductID:       productID,
		StartsAtUtc:     startAt,
		EndsAtUtc:       &end,
		Status:          "Scheduled",
		Notes:           marker,
		CreatedByUserID: actor,
	}
	app.CreatedAtUtc = startAt.Add(-6 * time.Hour)
	if err := tx.Create(&app).Error; err != nil {
		return fmt.Errorf("seed: create appointment %s: %w", marker, err)
	}
	return nil
}
