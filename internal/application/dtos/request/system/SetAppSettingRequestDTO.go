package system

type SetAppSettingRequestDTO struct {
	Value string `json:"value" binding:"required"`
}
