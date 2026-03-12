package customer

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/application/use_cases/common"
)

func (u *UseCase) SetActive(ctx context.Context, actor, id uuid.UUID, active bool) error {
	err := u.uow.Repositories().Customers().SetActive(ctx, id, active)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Customer", id.String(), map[bool]string{true: "ACTIVATE", false: "DEACTIVATE"}[active], nil)
	}
	return err
}
