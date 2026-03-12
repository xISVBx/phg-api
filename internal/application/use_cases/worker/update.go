package worker

import (
	"context"
	"time"

	"github.com/google/uuid"
	workerreq "photogallery/api_go/internal/application/dtos/request/worker"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) Update(ctx context.Context, actor, id uuid.UUID, in workerreq.UpdateWorkerRequestDTO) (*entities.Worker, error) {
	item, err := u.uow.Repositories().Workers().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	item.FullName, item.Phone, item.Email, item.FixedSalary, item.SalaryPeriod, item.Notes = in.FullName, in.Phone, in.Email, in.FixedSalary, in.SalaryPeriod, in.Notes
	item.UpdatedAtUtc = time.Now().UTC()
	err = u.uow.Repositories().Workers().Update(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Worker", item.ID.String(), "UPDATE", item)
	}
	return item, err
}
