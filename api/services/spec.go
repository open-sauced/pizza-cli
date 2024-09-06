package services

type MetaData struct {
	Page            int  `json:"page"`
	Limit           int  `json:"limit"`
	ItemCount       int  `json:"itemCount"`
	PageCount       int  `json:"pageCount"`
	HasPreviousPage bool `json:"hasPreviousPage"`
	HasNextPage     bool `json:"hasNextPage"`
}
