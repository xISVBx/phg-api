package workorder

type CreateWorkOrderRequestDTO struct {
	SaleID            string `json:"saleId" binding:"required,uuid"`
	Status            string `json:"status"`
	DueDateUtc        string `json:"dueDateUtc"`
	ResponsibleUserID string `json:"responsibleUserId"`
	Notes             string `json:"notes"`
}
