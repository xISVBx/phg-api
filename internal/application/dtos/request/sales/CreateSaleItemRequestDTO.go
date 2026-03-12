package sales

type CreateSaleItemRequestDTO struct {
	ProductID      string  `json:"productId" binding:"required,uuid"`
	Quantity       int     `json:"quantity" binding:"required,min=1"`
	UnitPrice      float64 `json:"unitPrice"`
	Discount       float64 `json:"discount"`
	DiscountReason string  `json:"discountReason"`
	Notes          string  `json:"notes"`
}
