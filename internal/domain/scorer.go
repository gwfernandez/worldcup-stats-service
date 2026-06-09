package domain

// Scorer represents a player in the historical scorers table.
type Scorer struct {
	FullName          string   `json:"fullName"`
	TeamCode          string   `json:"teamCode"`
	TeamName          string   `json:"teamName"`
	Goals             int32    `json:"goals"`
	ListTeams         []string `json:"listTeams"`
	ConfederationCode string   `json:"confederationCode"`
}

// ScorerFilter defines filters and pagination for listing historical scorers.
type ScorerFilter struct {
	Name              string
	Language          string
	TeamCode          string
	ConfederationCode string
	Page              int
	Size              int
}

// ScorerListResponse represents paginated historical scorers response.
type ScorerListResponse struct {
	Data       []Scorer       `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}
