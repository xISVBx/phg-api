package customer

type CreateCustomerRequestDTO struct {
	FullName     string `json:"fullName" binding:"required"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	CustomerCode string `json:"customerCode"`
	Document     string `json:"document"`
	Notes        string `json:"notes"`
}
