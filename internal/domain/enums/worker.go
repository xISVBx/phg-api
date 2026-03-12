package enums

type WorkerPaymentType string

type CommissionEntryStatus string

const (
	WorkerPaymentCommission WorkerPaymentType = "Commission"
	WorkerPaymentSalary     WorkerPaymentType = "Salary"

	CommissionEarned  CommissionEntryStatus = "Earned"
	CommissionPartial CommissionEntryStatus = "PartiallyPaid"
	CommissionPaid    CommissionEntryStatus = "Paid"
	CommissionVoid    CommissionEntryStatus = "Voided"
	CommissionAdjust  CommissionEntryStatus = "Adjusted"
)
