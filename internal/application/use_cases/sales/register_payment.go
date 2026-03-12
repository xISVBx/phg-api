package sales

import (
	"context"
	"time"

	"github.com/google/uuid"

	salesreq "photogallery/api_go/internal/application/dtos/request/sales"
	appif "photogallery/api_go/internal/application/interfaces"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/domain/enums"
)

func (u *UseCase) RegisterPayment(ctx context.Context, actor, saleID uuid.UUID, in salesreq.RegisterSalePaymentRequestDTO) (*entities.Sale, error) {
	now := time.Now().UTC()
	var sale *entities.Sale
	err := u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
		s, err := repos.Sales().GetByID(ctx, saleID)
		if err != nil {
			return err
		}
		payment := &entities.SalePayment{SaleID: saleID, Method: in.Method, Amount: in.Amount, Reference: in.Reference, PaidAtUtc: now, CreatedByUserID: actor}
		if err := repos.SalePayments().Create(ctx, payment); err != nil {
			return err
		}
		paid, err := repos.SalePayments().SumBySale(ctx, saleID)
		if err != nil {
			return err
		}
		if paid <= 0 {
			s.Status = string(enums.SaleStatusPending)
		} else if paid < s.Total {
			s.Status = string(enums.SaleStatusAbonada)
		} else {
			s.Status = string(enums.SaleStatusPagada)
		}
		if err := repos.Sales().Update(ctx, s); err != nil {
			return err
		}
		common.CreateAudit(ctx, repos, &actor, "SalePayment", payment.ID.String(), "REGISTER", payment)
		sale = s
		return nil
	})
	return sale, err
}
