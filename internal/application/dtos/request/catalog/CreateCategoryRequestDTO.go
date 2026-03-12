package catalog

type CreateCategoryRequestDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
