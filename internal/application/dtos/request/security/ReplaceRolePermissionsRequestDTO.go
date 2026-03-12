package security

type ReplaceRolePermissionsRequestDTO struct {
	Items []RolePermissionSetItemDTO `json:"items"`
}
