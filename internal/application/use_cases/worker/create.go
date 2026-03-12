package worker

import (
	"context"

	"github.com/google/uuid"
	workerreq "photogallery/api_go/internal/application/dtos/request/worker"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) Create(ctx context.Context, actor uuid.UUID, in workerreq.CreateWorkerRequestDTO) (*entities.Worker, error) {
	item := &entities.Worker{FullName: in.FullName, Phone: in.Phone, Email: in.Email, FixedSalary: in.FixedSalary, SalaryPeriod: in.SalaryPeriod, Notes: in.Notes, IsActive: true}
	err := u.uow.Repositories().Workers().Create(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Worker", item.ID.String(), "CREATE", item)
	}
	return item, err
}
