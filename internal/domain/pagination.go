package domain

// PaginationInfo represents metadata about a paginated response.
type PaginationInfo struct {
	Page          int   `json:"page"`
	Size          int   `json:"size"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
	HasNext       bool  `json:"hasNext"`
	HasPrevious   bool  `json:"hasPrevious"`
}
