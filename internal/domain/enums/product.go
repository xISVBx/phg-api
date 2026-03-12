package enums

type ProductType string

type CommissionType string

const (
	ProductPhysical ProductType = "Physical"
	ProductService  ProductType = "Service"
	ProductStudy    ProductType = "Study"

	CommissionNone    CommissionType = "None"
	CommissionPercent CommissionType = "Percent"
	CommissionFixed   CommissionType = "Fixed"
)
