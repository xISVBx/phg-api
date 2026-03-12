package enums

type WorkOrderStatus string

const (
	WorkOrderCreated   WorkOrderStatus = "Creado"
	WorkOrderInDev     WorkOrderStatus = "EnDesarrollo"
	WorkOrderReady     WorkOrderStatus = "Listo"
	WorkOrderDelivered WorkOrderStatus = "Entregado"
)
