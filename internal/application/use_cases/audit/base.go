package audit

import appif "photogallery/api_go/internal/application/interfaces"

type UseCase struct{ uow appif.IUnitOfWork }

func New(uow appif.IUnitOfWork) *UseCase { return &UseCase{uow: uow} }
