package security

type UpdateUserRequestDTO struct {
	FullName    string                     `json:"fullName" binding:"required"`
	Phone       string                     `json:"phone"`
	Email       string                     `json:"email"`
	IsActive    *bool                      `json:"isActive"`
	RoleID      string                     `json:"roleId" binding:"required,uuid"`
	Permissions []RolePermissionSetItemDTO `json:"permissions" binding:"required"`
}
