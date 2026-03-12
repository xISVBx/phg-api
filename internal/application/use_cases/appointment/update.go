package appointment

import (
	"context"
	"time"

	"github.com/google/uuid"
	appointmentreq "photogallery/api_go/internal/application/dtos/request/appointment"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/domain/enums"
)

func (u *UseCase) Update(ctx context.Context, actor, id uuid.UUID, in appointmentreq.UpdateAppointmentRequestDTO) (*entities.Appointment, error) {
	item, err := u.uow.Repositories().Appointments().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	cid, err := uuid.Parse(in.CustomerID)
	if err != nil {
		return nil, err
	}
	pid, err := uuid.Parse(in.ProductID)
	if err != nil {
		return nil, err
	}
	starts, err := time.Parse(time.RFC3339, in.StartsAtUtc)
	if err != nil {
		return nil, err
	}
	item.CustomerID, item.ProductID, item.StartsAtUtc, item.Status, item.Notes = cid, pid, starts, in.Status, in.Notes
	if item.Status == "" {
		item.Status = string(enums.AppointmentProgrammed)
	}
	if in.SaleID != nil && *in.SaleID != "" {
		sid, err := uuid.Parse(*in.SaleID)
		if err != nil {
			return nil, err
		}
		item.SaleID = &sid
	}
	if in.EndsAtUtc != "" {
		v, err := time.Parse(time.RFC3339, in.EndsAtUtc)
		if err != nil {
			return nil, err
		}
		item.EndsAtUtc = &v
	}
	item.UpdatedAtUtc = time.Now().UTC()
	err = u.uow.Repositories().Appointments().Update(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Appointment", item.ID.String(), "UPDATE", item)
	}
	return item, err
}
