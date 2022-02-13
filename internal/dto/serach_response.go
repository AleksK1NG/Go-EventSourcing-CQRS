package dto

type OrderSearchResponseDto struct {
	Pagination Pagination         `json:"pagination"`
	Orders     []OrderResponseDto `json:"orders"`
}

type Pagination struct {
	TotalCount int64 `json:"totalCount"`
	TotalPages int64 `json:"totalPages"`
	Page       int64 `json:"page"`
	Size       int64 `json:"size"`
	HasMore    bool  `json:"hasMore"`
}
