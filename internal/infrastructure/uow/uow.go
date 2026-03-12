package uow

import (
	"context"

	"gorm.io/gorm"

	appif "photogallery/api_go/internal/application/interfaces"
	drepo "photogallery/api_go/internal/domain/repositories"
	appointmentsrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/appointments"
	auditlogsrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/auditlogs"
	cashcategoriesrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/cashcategories"
	cashmovementsrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/cashmovements"
	categoriesrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/categories"
	commissionsrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/commissions"
	customersrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/customers"
	filelinksrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/filelinks"
	filesrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/files"
	menusrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/menus"
	overridesrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/overrides"
	permissionsrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/permissions"
	productsrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/products"
	rolepermissionsrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/rolepermissions"
	rolesrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/roles"
	saleitemsrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/saleitems"
	salepaymentsrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/salepayments"
	salesrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/sales"
	settingsrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/settings"
	submenusrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/submenus"
	userrolesrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/userroles"
	usersrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/users"
	workerpaymentsrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/workerpayments"
	workersrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/workers"
	workordersrepo "photogallery/api_go/internal/infrastructure/persistence/repositories/workorders"
)

type repositorySet struct {
	users          drepo.UserRepository
	roles          drepo.RoleRepository
	menus          drepo.MenuRepository
	subMenus       drepo.SubMenuRepository
	permissions    drepo.PermissionRepository
	userRoles      drepo.UserRoleRepository
	rolePerms      drepo.RolePermissionRepository
	overrides      drepo.OverrideRepository
	audit          drepo.AuditLogRepository
	categories     drepo.CategoryRepository
	products       drepo.ProductRepository
	customers      drepo.CustomerRepository
	sales          drepo.SaleRepository
	saleItems      drepo.SaleItemRepository
	salePayments   drepo.SalePaymentRepository
	workOrders     drepo.WorkOrderRepository
	appointments   drepo.AppointmentRepository
	files          drepo.FileRepository
	fileLinks      drepo.FileLinkRepository
	cashCategories drepo.CashCategoryRepository
	cashMovements  drepo.CashMovementRepository
	workers        drepo.WorkerRepository
	commissions    drepo.CommissionRepository
	workerPayments drepo.WorkerPaymentRepository
	settings       drepo.AppSettingRepository
}

func newRepositorySet(db *gorm.DB) *repositorySet {
	return &repositorySet{
		users: usersrepo.NewRepository(db), roles: rolesrepo.NewRepository(db), menus: menusrepo.NewRepository(db),
		subMenus: submenusrepo.NewRepository(db), permissions: permissionsrepo.NewRepository(db), userRoles: userrolesrepo.NewRepository(db),
		rolePerms: rolepermissionsrepo.NewRepository(db), overrides: overridesrepo.NewRepository(db), audit: auditlogsrepo.NewRepository(db),
		categories: categoriesrepo.NewRepository(db), products: productsrepo.NewRepository(db), customers: customersrepo.NewRepository(db),
		sales: salesrepo.NewRepository(db), saleItems: saleitemsrepo.NewRepository(db), salePayments: salepaymentsrepo.NewRepository(db),
		workOrders: workordersrepo.NewRepository(db), appointments: appointmentsrepo.NewRepository(db), files: filesrepo.NewRepository(db),
		fileLinks: filelinksrepo.NewRepository(db), cashCategories: cashcategoriesrepo.NewRepository(db), cashMovements: cashmovementsrepo.NewRepository(db),
		workers: workersrepo.NewRepository(db), commissions: commissionsrepo.NewRepository(db), workerPayments: workerpaymentsrepo.NewRepository(db),
		settings: settingsrepo.NewRepository(db),
	}
}

func (r *repositorySet) Users() drepo.UserRepository                     { return r.users }
func (r *repositorySet) Roles() drepo.RoleRepository                     { return r.roles }
func (r *repositorySet) Menus() drepo.MenuRepository                     { return r.menus }
func (r *repositorySet) SubMenus() drepo.SubMenuRepository               { return r.subMenus }
func (r *repositorySet) Permissions() drepo.PermissionRepository         { return r.permissions }
func (r *repositorySet) UserRoles() drepo.UserRoleRepository             { return r.userRoles }
func (r *repositorySet) RolePermissions() drepo.RolePermissionRepository { return r.rolePerms }
func (r *repositorySet) Overrides() drepo.OverrideRepository             { return r.overrides }
func (r *repositorySet) AuditLogs() drepo.AuditLogRepository             { return r.audit }
func (r *repositorySet) Categories() drepo.CategoryRepository            { return r.categories }
func (r *repositorySet) Products() drepo.ProductRepository               { return r.products }
func (r *repositorySet) Customers() drepo.CustomerRepository             { return r.customers }
func (r *repositorySet) Sales() drepo.SaleRepository                     { return r.sales }
func (r *repositorySet) SaleItems() drepo.SaleItemRepository             { return r.saleItems }
func (r *repositorySet) SalePayments() drepo.SalePaymentRepository       { return r.salePayments }
func (r *repositorySet) WorkOrders() drepo.WorkOrderRepository           { return r.workOrders }
func (r *repositorySet) Appointments() drepo.AppointmentRepository       { return r.appointments }
func (r *repositorySet) Files() drepo.FileRepository                     { return r.files }
func (r *repositorySet) FileLinks() drepo.FileLinkRepository             { return r.fileLinks }
func (r *repositorySet) CashCategories() drepo.CashCategoryRepository    { return r.cashCategories }
func (r *repositorySet) CashMovements() drepo.CashMovementRepository     { return r.cashMovements }
func (r *repositorySet) Workers() drepo.WorkerRepository                 { return r.workers }
func (r *repositorySet) Commissions() drepo.CommissionRepository         { return r.commissions }
func (r *repositorySet) WorkerPayments() drepo.WorkerPaymentRepository   { return r.workerPayments }
func (r *repositorySet) Settings() drepo.AppSettingRepository            { return r.settings }

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork { return &UnitOfWork{db: db} }

func (u *UnitOfWork) Repositories() appif.RepositorySet { return newRepositorySet(u.db) }

func (u *UnitOfWork) Transaction(ctx context.Context, fn func(repos appif.RepositorySet) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(newRepositorySet(tx))
	})
}
