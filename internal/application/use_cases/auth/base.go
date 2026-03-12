package auth

import (
	appif "photogallery/api_go/internal/application/interfaces"
	appsvc "photogallery/api_go/internal/application/interfaces/services"
)

type UseCase struct {
	uow  appif.IUnitOfWork
	jwt  appsvc.IJWTService
	pass appsvc.IPasswordService
}

func New(uow appif.IUnitOfWork, jwt appsvc.IJWTService, pass appsvc.IPasswordService) *UseCase {
	return &UseCase{uow: uow, jwt: jwt, pass: pass}
}
