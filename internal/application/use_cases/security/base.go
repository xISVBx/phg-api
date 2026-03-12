package security

import (
	appif "photogallery/api_go/internal/application/interfaces"
	appsvc "photogallery/api_go/internal/application/interfaces/services"
)

type UseCase struct {
	uow  appif.IUnitOfWork
	pass appsvc.IPasswordService
}

func New(uow appif.IUnitOfWork, pass appsvc.IPasswordService) *UseCase {
	return &UseCase{uow: uow, pass: pass}
}
