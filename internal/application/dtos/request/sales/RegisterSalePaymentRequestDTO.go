package sales

type RegisterSalePaymentRequestDTO struct {
	Method    string  `json:"method" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Reference string  `json:"reference"`
}
