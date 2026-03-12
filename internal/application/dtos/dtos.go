package dtos

type AuthLoginDTO struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	ExpiresIn    int64       `json:"expiresIn"`
	User         interface{} `json:"user"`
}

type EffectivePermissionNode struct {
	MenuCode string             `json:"menuCode"`
	MenuName string             `json:"menuName"`
	SubMenus []EffectiveSubMenu `json:"subMenus"`
}

type EffectiveSubMenu struct {
	SubMenuCode string   `json:"subMenuCode"`
	SubMenuName string   `json:"subMenuName"`
	Permissions []string `json:"permissions"`
}
