package worker

import (
	"context"
	"time"

	"github.com/google/uuid"

	workerreq "photogallery/api_go/internal/application/dtos/request/worker"
	appif "photogallery/api_go/internal/application/interfaces"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/domain/enums"
)

func (u *UseCase) PaySalary(ctx context.Context, actor uuid.UUID, in workerreq.PaySalaryRequestDTO) error {
	wid, err := uuid.Parse(in.WorkerID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	return u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
		wp := &entities.WorkerPayment{WorkerID: wid, Type: string(enums.WorkerPaymentSalary), Method: in.Method, Amount: in.Amount, Notes: in.Notes, PaidAtUtc: now, CreatedByUserID: actor}
		if err := repos.WorkerPayments().Create(ctx, wp); err != nil {
			return err
		}
		cm := &entities.CashMovement{Type: string(enums.CashOut), Method: in.Method, Amount: in.Amount, RelatedEntityType: "WorkerPayment", RelatedEntityID: wp.ID.String(), Notes: "Pago sueldo", CreatedByUserID: actor}
		cm.CategoryID = uuid.Nil
		if err := repos.CashMovements().Create(ctx, cm); err != nil {
			return err
		}
		common.CreateAudit(ctx, repos, &actor, "WorkerPayment", wp.ID.String(), "PAY_SALARY", wp)
		return nil
	})
}
