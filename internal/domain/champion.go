package domain

// Champion represents a team in the historical champions table.
type Champion struct {
	Team              SimpleTeam `json:"team"`
	Wins              int64      `json:"wins"`
	Years             []int32    `json:"years"`
	ConfederationCode string     `json:"confederationCode"`
}

// ChampionFilter defines pagination for listing champions.
type ChampionFilter struct {
	Language string
	Page     int
	Size     int
}

// ChampionListResponse represents paginated champions response.
type ChampionListResponse struct {
	Data       []Champion     `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}
