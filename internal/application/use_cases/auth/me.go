package auth

import (
	"context"

	"github.com/google/uuid"
	authresp "photogallery/api_go/internal/application/dtos/responses/auth"
)

func (u *UseCase) Me(ctx context.Context, userID uuid.UUID) (*authresp.MeResponseDTO, error) {
	repos := u.uow.Repositories()

	user, err := repos.Users().GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userRoles, err := repos.UserRoles().ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	roles := make([]authresp.MeRoleDTO, 0, len(userRoles))
	for _, ur := range userRoles {
		role, err := repos.Roles().GetByID(ctx, ur.RoleID)
		if err != nil {
			return nil, err
		}
		roles = append(roles, authresp.MeRoleDTO{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			IsActive:    role.IsActive,
			IsPrimary:   ur.IsPrimary,
			RoleType:    role.RoleType,
		})
	}

	perms, err := u.MyPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &authresp.MeResponseDTO{
		ID:           user.ID,
		Username:     user.Username,
		FullName:     user.FullName,
		Phone:        user.Phone,
		Email:        user.Email,
		IsActive:     user.IsActive,
		CreatedAtUtc: user.CreatedAtUtc,
		UpdatedAtUtc: user.UpdatedAtUtc,
		Roles:        roles,
		Permissions:  perms,
	}, nil
}
