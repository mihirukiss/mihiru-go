package models

type PageParams struct {
	PageSize  *int64 `json:"pageSize"`
	PageIndex *int64 `json:"pageIndex"`
}

type PageResult struct {
	PageParams
	PageCount int64 `json:"pageCount"`
	Count     int64 `json:"count"`
}
