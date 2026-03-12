package appointment

type CreateAppointmentRequestDTO struct {
	CustomerID  string  `json:"customerId" binding:"required,uuid"`
	SaleID      *string `json:"saleId"`
	ProductID   string  `json:"productId" binding:"required,uuid"`
	StartsAtUtc string  `json:"startsAtUtc" binding:"required"`
	EndsAtUtc   string  `json:"endsAtUtc"`
	Status      string  `json:"status"`
	Notes       string  `json:"notes"`
}
