package security

type CreateUserRequestDTO struct {
	Username    string                     `json:"username" binding:"required"`
	Password    string                     `json:"password" binding:"required,min=6"`
	FullName    string                     `json:"fullName" binding:"required"`
	Phone       string                     `json:"phone"`
	Email       string                     `json:"email"`
	RoleID      string                     `json:"roleId" binding:"required,uuid"`
	Permissions []RolePermissionSetItemDTO `json:"permissions" binding:"required"`
}
