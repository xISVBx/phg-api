package cash

type CreateCashMovementRequestDTO struct {
	Type              string  `json:"type" binding:"required"`
	CategoryID        string  `json:"categoryId" binding:"required,uuid"`
	Method            string  `json:"method" binding:"required"`
	Amount            float64 `json:"amount" binding:"required,gt=0"`
	Reference         string  `json:"reference"`
	RelatedEntityType string  `json:"relatedEntityType"`
	RelatedEntityID   string  `json:"relatedEntityId"`
	Notes             string  `json:"notes"`
}
