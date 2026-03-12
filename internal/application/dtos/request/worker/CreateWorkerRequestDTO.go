package worker

type CreateWorkerRequestDTO struct {
	FullName     string  `json:"fullName" binding:"required"`
	Phone        string  `json:"phone"`
	Email        string  `json:"email"`
	FixedSalary  float64 `json:"fixedSalary"`
	SalaryPeriod string  `json:"salaryPeriod"`
	Notes        string  `json:"notes"`
}
