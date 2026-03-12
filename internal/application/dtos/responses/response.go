package responses

type PageMeta struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

type PageResponse[T any] struct {
	Items []T      `json:"items"`
	Meta  PageMeta `json:"meta"`
}
