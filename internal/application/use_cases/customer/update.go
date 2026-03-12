package customer

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	customerreq "photogallery/api_go/internal/application/dtos/request/customer"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) Update(ctx context.Context, actor, id uuid.UUID, in customerreq.UpdateCustomerRequestDTO) (*entities.Customer, error) {
	item, err := u.uow.Repositories().Customers().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.CustomerCode == "" {
		in.CustomerCode = item.CustomerCode
	}
	item.FullName, item.Phone, item.Email, item.CustomerCode, item.Document, item.Notes = in.FullName, in.Phone, strings.TrimSpace(strings.ToLower(in.Email)), in.CustomerCode, in.Document, in.Notes
	item.UpdatedAtUtc = time.Now().UTC()
	err = u.uow.Repositories().Customers().Update(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Customer", item.ID.String(), "UPDATE", item)
	}
	return item, err
}
