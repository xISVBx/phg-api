package worker

type PaySalaryRequestDTO struct {
	WorkerID string  `json:"workerId" binding:"required,uuid"`
	Method   string  `json:"method" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Notes    string  `json:"notes"`
}
