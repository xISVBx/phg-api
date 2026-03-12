package mappers

import (
	securityresp "photogallery/api_go/internal/application/dtos/responses/security"
	"photogallery/api_go/internal/domain/entities"
)

func ToUserResponse(u *entities.User) securityresp.UserResponseDTO {
	return securityresp.UserResponseDTO{
		ID:           u.ID,
		Username:     u.Username,
		FullName:     u.FullName,
		Phone:        u.Phone,
		Email:        u.Email,
		IsActive:     u.IsActive,
		CreatedAtUtc: u.CreatedAtUtc,
		UpdatedAtUtc: u.UpdatedAtUtc,
	}
}

func ToUserResponseList(users []entities.User) []securityresp.UserResponseDTO {
	out := make([]securityresp.UserResponseDTO, 0, len(users))
	for i := range users {
		item := users[i]
		out = append(out, ToUserResponse(&item))
	}
	return out
}
