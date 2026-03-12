package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	authreq "photogallery/api_go/internal/application/dtos/request/auth"
)

func (u *UseCase) ChangePassword(ctx context.Context, userID uuid.UUID, in authreq.ChangePasswordRequestDTO) error {
	user, err := u.uow.Repositories().Users().GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !u.pass.Compare(user.PasswordHash, in.OldPassword) {
		return errors.New("invalid old password")
	}
	h, err := u.pass.Hash(in.NewPassword)
	if err != nil {
		return err
	}
	return u.uow.Repositories().Users().SetPassword(ctx, userID, h)
}
