package enums

type SaleStatus string

const (
	SaleStatusPending SaleStatus = "Pending"
	SaleStatusAbonada SaleStatus = "Abonada"
	SaleStatusPagada  SaleStatus = "Pagada"
	SaleStatusAnulada SaleStatus = "Anulada"
	SaleStatusCancel  SaleStatus = "Cancelada"
)
