package use_cases

import (
	appif "photogallery/api_go/internal/application/interfaces"
	appsvc "photogallery/api_go/internal/application/interfaces/services"
	"photogallery/api_go/internal/application/use_cases/appointment"
	"photogallery/api_go/internal/application/use_cases/audit"
	"photogallery/api_go/internal/application/use_cases/auth"
	"photogallery/api_go/internal/application/use_cases/cash"
	"photogallery/api_go/internal/application/use_cases/catalog"
	"photogallery/api_go/internal/application/use_cases/customer"
	"photogallery/api_go/internal/application/use_cases/files"
	"photogallery/api_go/internal/application/use_cases/sales"
	"photogallery/api_go/internal/application/use_cases/security"
	"photogallery/api_go/internal/application/use_cases/system"
	"photogallery/api_go/internal/application/use_cases/worker"
	"photogallery/api_go/internal/application/use_cases/workorder"
)

type UseCases struct {
	Auth         *auth.UseCase
	Security     *security.UseCase
	Catalog      *catalog.UseCase
	Customers    *customer.UseCase
	Sales        *sales.UseCase
	WorkOrders   *workorder.UseCase
	Appointments *appointment.UseCase
	Files        *files.UseCase
	Cash         *cash.UseCase
	Workers      *worker.UseCase
	Audit        *audit.UseCase
	System       *system.UseCase
}

func NewUseCases(uow appif.IUnitOfWork, jwt appsvc.IJWTService, pass appsvc.IPasswordService, fs appsvc.IFileStorageService, notifier appsvc.INotificationService) *UseCases {
	return &UseCases{
		Auth:         auth.New(uow, jwt, pass),
		Security:     security.New(uow, pass),
		Catalog:      catalog.New(uow),
		Customers:    customer.New(uow),
		Sales:        sales.New(uow, notifier),
		WorkOrders:   workorder.New(uow),
		Appointments: appointment.New(uow),
		Files:        files.New(uow, fs),
		Cash:         cash.New(uow),
		Workers:      worker.New(uow),
		Audit:        audit.New(uow),
		System:       system.New(uow),
	}
}
