package queries

type PageQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Q        string `form:"q"`
	Sort     string `form:"sort"`
	Dir      string `form:"dir"`
}
