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

func (u *UseCase) PayCommissionFIFO(ctx context.Context, actor uuid.UUID, in workerreq.PayCommissionRequestDTO) error {
	wid, err := uuid.Parse(in.WorkerID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	return u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
		pending, err := repos.Commissions().ListPendingByWorker(ctx, wid)
		if err != nil {
			return err
		}
		left := in.Amount
		allocs := make([]entities.WorkerPaymentAllocation, 0)
		payment := &entities.WorkerPayment{WorkerID: wid, Type: string(enums.WorkerPaymentCommission), Method: in.Method, Amount: in.Amount, Notes: in.Notes, PaidAtUtc: now, CreatedByUserID: actor}
		if err := repos.WorkerPayments().Create(ctx, payment); err != nil {
			return err
		}
		for _, ce := range pending {
			if left <= 0 {
				break
			}
			pendingAmount := ce.Amount - ce.PaidAmount
			if pendingAmount <= 0 {
				continue
			}
			apply := pendingAmount
			if left < apply {
				apply = left
			}
			ce.PaidAmount += apply
			if ce.PaidAmount >= ce.Amount {
				ce.Status = string(enums.CommissionPaid)
			} else {
				ce.Status = string(enums.CommissionPartial)
			}
			if err := repos.Commissions().Update(ctx, &ce); err != nil {
				return err
			}
			allocs = append(allocs, entities.WorkerPaymentAllocation{WorkerPaymentID: payment.ID, CommissionEntryID: ce.ID, AmountApplied: apply})
			left -= apply
		}
		if err := repos.WorkerPayments().CreateAllocations(ctx, allocs); err != nil {
			return err
		}
		cm := &entities.CashMovement{Type: string(enums.CashOut), Method: in.Method, Amount: in.Amount, RelatedEntityType: "WorkerPayment", RelatedEntityID: payment.ID.String(), Notes: "Pago comisiones", CreatedByUserID: actor}
		cm.CategoryID = uuid.Nil
		if err := repos.CashMovements().Create(ctx, cm); err != nil {
			return err
		}
		// NOTE: no se recrea el WorkerPayment para evitar violar PK duplicada.
		// Si se requiere persistir CashMovementID, agregar Update al repo y usarlo aquí.
		common.CreateAudit(ctx, repos, &actor, "WorkerPayment", payment.ID.String(), "PAY_COMMISSIONS", map[string]any{"requested": in.Amount, "remaining": left})
		return nil
	})
}
