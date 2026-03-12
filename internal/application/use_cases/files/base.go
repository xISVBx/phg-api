package files

import (
	appif "photogallery/api_go/internal/application/interfaces"
	appsvc "photogallery/api_go/internal/application/interfaces/services"
)

type UseCase struct {
	uow     appif.IUnitOfWork
	storage appsvc.IFileStorageService
}

func New(uow appif.IUnitOfWork, storage appsvc.IFileStorageService) *UseCase {
	return &UseCase{uow: uow, storage: storage}
}
