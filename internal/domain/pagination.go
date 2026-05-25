package domain

// PaginationInfo represents metadata about a paginated response.
type PaginationInfo struct {
	Page          int   `json:"page"`
	Size          int   `json:"size"`
	TotalElements int64 `json:"total_elements"`
	TotalPages    int   `json:"total_pages"`
	HasNext       bool  `json:"has_next"`
	HasPrevious   bool  `json:"has_previous"`
}
