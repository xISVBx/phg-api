package auth

import (
	"context"

	"photogallery/api_go/internal/application/dtos"
	authreq "photogallery/api_go/internal/application/dtos/request/auth"
)

func (u *UseCase) Refresh(ctx context.Context, in authreq.RefreshRequestDTO) (*dtos.AuthLoginDTO, error) {
	claims, err := u.jwt.Parse(in.RefreshToken)
	if err != nil {
		return nil, err
	}
	user, err := u.uow.Repositories().Users().GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	tokens, err := u.jwt.Generate(user.ID, user.Username)
	if err != nil {
		return nil, err
	}
	return &dtos.AuthLoginDTO{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken, ExpiresIn: tokens.ExpiresIn, User: map[string]any{"id": user.ID, "fullName": user.FullName}}, nil
}
