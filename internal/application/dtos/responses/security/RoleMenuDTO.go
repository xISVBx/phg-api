package security

import "github.com/google/uuid"

type RoleMenuDTO struct {
	ID       uuid.UUID        `json:"id"`
	Code     string           `json:"code"`
	Name     string           `json:"name"`
	SubMenus []RoleSubMenuDTO `json:"subMenus"`
}
