package repositories

type QueryOptions struct {
	Page     int
	PageSize int
	Q        string
	Sort     string
	Dir      string
}
