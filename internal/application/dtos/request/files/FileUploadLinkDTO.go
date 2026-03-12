package files

type FileUploadLinkDTO struct {
	EntityType string `form:"entityType" binding:"required"`
	EntityID   string `form:"entityId" binding:"required,uuid"`
	CustomerID string `form:"customerId"`
	Notes      string `form:"notes"`
}
