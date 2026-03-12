package security

type RolePermissionSetItemDTO struct {
	MenuID        string   `json:"menuId" binding:"required,uuid"`
	SubMenuID     string   `json:"subMenuId" binding:"required,uuid"`
	PermissionIDs []string `json:"permissionIds" binding:"required"`
}
