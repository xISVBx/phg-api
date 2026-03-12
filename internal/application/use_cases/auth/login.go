package auth

import (
	"context"
	"errors"

	"photogallery/api_go/internal/application/dtos"
	authreq "photogallery/api_go/internal/application/dtos/request/auth"
)

func (u *UseCase) Login(ctx context.Context, in authreq.LoginRequestDTO) (*dtos.AuthLoginDTO, error) {
	user, err := u.uow.Repositories().Users().GetByUsernameOrEmail(ctx, in.Username)
	if err != nil || !user.IsActive {
		return nil, errors.New("invalid credentials")
	}
	if !u.pass.Compare(user.PasswordHash, in.Password) {
		return nil, errors.New("invalid credentials")
	}
	tokens, err := u.jwt.Generate(user.ID, user.Username)
	if err != nil {
		return nil, err
	}
	return &dtos.AuthLoginDTO{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken, ExpiresIn: tokens.ExpiresIn, User: map[string]any{"id": user.ID, "fullName": user.FullName}}, nil
}
