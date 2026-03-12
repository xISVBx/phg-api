package interfaces

import (
	"context"

	"photogallery/api_go/internal/domain/repositories"
)

type RepositorySet interface {
	Users() repositories.UserRepository
	Roles() repositories.RoleRepository
	Menus() repositories.MenuRepository
	SubMenus() repositories.SubMenuRepository
	Permissions() repositories.PermissionRepository
	UserRoles() repositories.UserRoleRepository
	RolePermissions() repositories.RolePermissionRepository
	Overrides() repositories.OverrideRepository
	AuditLogs() repositories.AuditLogRepository
	Categories() repositories.CategoryRepository
	Products() repositories.ProductRepository
	Customers() repositories.CustomerRepository
	Sales() repositories.SaleRepository
	SaleItems() repositories.SaleItemRepository
	SalePayments() repositories.SalePaymentRepository
	WorkOrders() repositories.WorkOrderRepository
	Appointments() repositories.AppointmentRepository
	Files() repositories.FileRepository
	FileLinks() repositories.FileLinkRepository
	CashCategories() repositories.CashCategoryRepository
	CashMovements() repositories.CashMovementRepository
	Workers() repositories.WorkerRepository
	Commissions() repositories.CommissionRepository
	WorkerPayments() repositories.WorkerPaymentRepository
	Settings() repositories.AppSettingRepository
}

type IUnitOfWork interface {
	Repositories() RepositorySet
	Transaction(ctx context.Context, fn func(repos RepositorySet) error) error
}
