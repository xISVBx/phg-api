package workorder

import (
	"context"
	"time"

	"github.com/google/uuid"
	workorderreq "photogallery/api_go/internal/application/dtos/request/workorder"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/domain/enums"
)

func (u *UseCase) Create(ctx context.Context, actor uuid.UUID, in workorderreq.CreateWorkOrderRequestDTO) (*entities.WorkOrder, error) {
	sid, err := uuid.Parse(in.SaleID)
	if err != nil {
		return nil, err
	}
	item := &entities.WorkOrder{SaleID: sid, Status: in.Status, Notes: in.Notes}
	if item.Status == "" {
		item.Status = string(enums.WorkOrderCreated)
	}
	if in.DueDateUtc != "" {
		d, err := time.Parse(time.RFC3339, in.DueDateUtc)
		if err != nil {
			return nil, err
		}
		item.DueDateUtc = &d
	}
	if in.ResponsibleUserID != "" {
		uid, err := uuid.Parse(in.ResponsibleUserID)
		if err != nil {
			return nil, err
		}
		item.ResponsibleUserID = &uid
	}
	err = u.uow.Repositories().WorkOrders().Create(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "WorkOrder", item.ID.String(), "CREATE", item)
	}
	return item, err
}
