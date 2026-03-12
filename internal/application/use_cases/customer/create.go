package customer

import (
	"context"
	"strings"

	"github.com/google/uuid"

	customerreq "photogallery/api_go/internal/application/dtos/request/customer"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) Create(ctx context.Context, actor uuid.UUID, in customerreq.CreateCustomerRequestDTO) (*entities.Customer, error) {
	code := strings.TrimSpace(in.CustomerCode)
	if code == "" {
		code = "CUST-" + strings.ToUpper(uuid.NewString()[:8])
	}
	item := &entities.Customer{FullName: in.FullName, Phone: in.Phone, Email: strings.TrimSpace(strings.ToLower(in.Email)), CustomerCode: code, Document: in.Document, Notes: in.Notes, IsActive: true}
	err := u.uow.Repositories().Customers().Create(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Customer", item.ID.String(), "CREATE", item)
	}
	return item, err
}
