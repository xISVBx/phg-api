package sales

import (
	appif "photogallery/api_go/internal/application/interfaces"
	appsvc "photogallery/api_go/internal/application/interfaces/services"
)

type UseCase struct {
	uow      appif.IUnitOfWork
	notifier appsvc.INotificationService
}

func New(uow appif.IUnitOfWork, notifier appsvc.INotificationService) *UseCase {
	return &UseCase{uow: uow, notifier: notifier}
}
