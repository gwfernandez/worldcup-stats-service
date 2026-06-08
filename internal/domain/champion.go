package domain

// Champion represents a team in the historical champions table.
type Champion struct {
	TeamCode string  `json:"teamCode"`
	Name     string  `json:"name"`
	Wins     int64   `json:"wins"`
	Years    []int32 `json:"years"`
}

// ChampionFilter defines pagination for listing champions.
type ChampionFilter struct {
	Page int
	Size int
}

// ChampionListResponse represents paginated champions response.
type ChampionListResponse struct {
	Data       []Champion     `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}
