package catalog

import (
	"context"

	"github.com/google/uuid"
	"photogallery/api_go/internal/application/use_cases/common"
)

func (u *UseCase) SetCategoryActive(ctx context.Context, actor, id uuid.UUID, active bool) error {
	err := u.uow.Repositories().Categories().SetActive(ctx, id, active)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "Category", id.String(), map[bool]string{true: "ACTIVATE", false: "DEACTIVATE"}[active], nil)
	}
	return err
}
