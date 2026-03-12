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

func (u *UseCase) Create(ctx context.Context, actor uuid.UUID, in appointmentreq.CreateAppointmentRequestDTO) (*entities.Appointment, error) {
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
	item := &entities.Appointment{CustomerID: cid, ProductID: pid, StartsAtUtc: starts, Status: in.Status, Notes: in.Notes, CreatedByUserID: actor}
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
	err = u.uow.Repositories().Appointments().Create(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Appointment", item.ID.String(), "CREATE", item)
	}
	return item, err
}
