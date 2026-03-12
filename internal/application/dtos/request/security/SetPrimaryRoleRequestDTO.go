package security

type SetPrimaryRoleRequestDTO struct {
	PrimaryRoleID string `json:"primaryRoleId" binding:"required,uuid"`
}
