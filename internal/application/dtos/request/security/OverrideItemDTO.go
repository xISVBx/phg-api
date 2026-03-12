package security

type OverrideItemDTO struct {
	SubMenuID    string `json:"subMenuId" binding:"required,uuid"`
	PermissionID string `json:"permissionId" binding:"required,uuid"`
	Mode         string `json:"mode" binding:"required,oneof=Grant Revoke"`
}
