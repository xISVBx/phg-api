package security

type SetPasswordRequestDTO struct {
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}
