package sales

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	salesreq "photogallery/api_go/internal/application/dtos/request/sales"
	appif "photogallery/api_go/internal/application/interfaces"
	appsvc "photogallery/api_go/internal/application/interfaces/services"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/domain/enums"
)

func (u *UseCase) Create(ctx context.Context, actor uuid.UUID, in salesreq.CreateSaleRequestDTO) (*entities.Sale, error) {
	now := time.Now().UTC()
	var out *entities.Sale
	err := u.uow.Transaction(ctx, func(repos appif.RepositorySet) error {
		var cid *uuid.UUID
		if in.CustomerID != nil && *in.CustomerID != "" {
			v, err := uuid.Parse(*in.CustomerID)
			if err != nil {
				return err
			}
			cid = &v
		}
		sale := &entities.Sale{CustomerID: cid, SellerUserID: actor, Status: string(enums.SaleStatusPending), NotifyOptIn: in.NotifyOptIn}
		var subtotal, discount, totalCost, totalComm float64
		saleItems := make([]entities.SaleItem, 0, len(in.Items))
		requiresDelivery := false
		for _, it := range in.Items {
			pid, err := uuid.Parse(it.ProductID)
			if err != nil {
				return err
			}
			prod, err := repos.Products().GetByID(ctx, pid)
			if err != nil {
				return err
			}
			unitPrice := it.UnitPrice
			if unitPrice == 0 {
				unitPrice = prod.BasePrice
			}
			line := unitPrice * float64(it.Quantity)
			subtotal += line
			discount += it.Discount
			totalCost += prod.Cost * float64(it.Quantity)
			if strings.EqualFold(prod.CommissionType, string(enums.CommissionPercent)) {
				totalComm += (unitPrice * (prod.CommissionValue / 100)) * float64(it.Quantity)
			}
			if strings.EqualFold(prod.CommissionType, string(enums.CommissionFixed)) {
				totalComm += prod.CommissionValue * float64(it.Quantity)
			}
			saleItems = append(saleItems, entities.SaleItem{ProductID: pid, Quantity: it.Quantity, UnitPriceSnapshot: unitPrice, UnitCostSnapshot: prod.Cost, CommissionTypeSnapshot: prod.CommissionType, CommissionValueSnapshot: prod.CommissionValue, DiscountSnapshot: it.Discount, DiscountReason: it.DiscountReason, Notes: it.Notes, RequiresDeliverySnapshot: prod.RequiresDelivery, LeadDaysSnapshot: prod.DefaultLeadDays})
			requiresDelivery = requiresDelivery || prod.RequiresDelivery
		}
		sale.Subtotal, sale.DiscountTotal, sale.Total = subtotal, discount, subtotal-discount
		sale.TotalCostSnapshot, sale.TotalCommissionSnapshot = totalCost, totalComm
		if err := repos.Sales().Create(ctx, sale); err != nil {
			return err
		}
		for i := range saleItems {
			saleItems[i].SaleID = sale.ID
		}
		if err := repos.SaleItems().CreateMany(ctx, saleItems); err != nil {
			return err
		}
		if requiresDelivery {
			due := now.AddDate(0, 0, 3)
			wo := &entities.WorkOrder{SaleID: sale.ID, Status: string(enums.WorkOrderCreated), DueDateUtc: &due}
			if err := repos.WorkOrders().Create(ctx, wo); err != nil {
				return err
			}
			common.CreateAudit(ctx, repos, &actor, "WorkOrder", wo.ID.String(), "AUTO_CREATE", wo)
		}
		common.CreateAudit(ctx, repos, &actor, "Sale", sale.ID.String(), "CREATE", sale)
		out = sale
		return nil
	})
	if err != nil {
		return nil, err
	}
	if out != nil && out.NotifyOptIn {
		_ = u.notifier.Send(ctx, appsvc.NotificationPayload{Channel: "none", Title: "TODO", Body: "TODO: integrar notificaciones reales"})
	}
	return out, nil
}
