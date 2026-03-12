package cash

import (
	"context"
	"time"

	"github.com/google/uuid"
	cashreq "photogallery/api_go/internal/application/dtos/request/cash"
	"photogallery/api_go/internal/application/use_cases/common"
	"photogallery/api_go/internal/domain/entities"
)

func (u *UseCase) UpdateMovement(ctx context.Context, actor, id uuid.UUID, in cashreq.UpdateCashMovementRequestDTO) (*entities.CashMovement, error) {
	item, err := u.uow.Repositories().CashMovements().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	cid, err := uuid.Parse(in.CategoryID)
	if err != nil {
		return nil, err
	}
	item.Type, item.CategoryID, item.Method, item.Amount, item.Reference, item.RelatedEntityType, item.RelatedEntityID, item.Notes = in.Type, cid, in.Method, in.Amount, in.Reference, in.RelatedEntityType, in.RelatedEntityID, in.Notes
	item.UpdatedAtUtc = time.Now().UTC()
	err = u.uow.Repositories().CashMovements().Update(ctx, item)
	if err == nil {
		common.CreateAudit(ctx, u.uow.Repositories(), &actor, "CashMovement", item.ID.String(), "UPDATE", item)
	}
	return item, err
}
