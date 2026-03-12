package cash

import (
	"context"

	"github.com/google/uuid"
	cashreq "photogallery/api_go/internal/application/dtos/request/cash"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) CreateMovement(ctx context.Context, actor uuid.UUID, in cashreq.CreateCashMovementRequestDTO) (*entities.CashMovement, error) {
	cid, err := uuid.Parse(in.CategoryID)
	if err != nil {
		return nil, err
	}
	item := &entities.CashMovement{Type: in.Type, CategoryID: cid, Method: in.Method, Amount: in.Amount, Reference: in.Reference, RelatedEntityType: in.RelatedEntityType, RelatedEntityID: in.RelatedEntityID, Notes: in.Notes, CreatedByUserID: actor}
	err = u.uow.Repositories().CashMovements().Create(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "CashMovement", item.ID.String(), "CREATE", item)
	}
	return item, err
}
