package enums

type CashType string

type CashCategoryType string

const (
	CashIn  CashType = "In"
	CashOut CashType = "Out"

	CashCategoryIn   CashCategoryType = "In"
	CashCategoryOut  CashCategoryType = "Out"
	CashCategoryBoth CashCategoryType = "Both"
)
