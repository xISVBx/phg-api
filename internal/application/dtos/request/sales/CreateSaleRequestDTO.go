package sales

type CreateSaleRequestDTO struct {
	CustomerID  *string                    `json:"customerId"`
	NotifyOptIn bool                       `json:"notifyOptIn"`
	Items       []CreateSaleItemRequestDTO `json:"items" binding:"required,min=1"`
}
