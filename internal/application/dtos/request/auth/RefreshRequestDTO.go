package auth

type RefreshRequestDTO struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
