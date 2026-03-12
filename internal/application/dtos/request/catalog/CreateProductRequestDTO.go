package catalog

type CreateProductRequestDTO struct {
	CategoryID       string  `json:"categoryId" binding:"required,uuid"`
	Name             string  `json:"name" binding:"required"`
	Type             string  `json:"type" binding:"required"`
	BasePrice        float64 `json:"basePrice"`
	Cost             float64 `json:"cost"`
	CommissionType   string  `json:"commissionType"`
	CommissionValue  float64 `json:"commissionValue"`
	RequiresDelivery bool    `json:"requiresDelivery"`
	DefaultLeadDays  *int    `json:"defaultLeadDays"`
	Notes            string  `json:"notes"`
}
