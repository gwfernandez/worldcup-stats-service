package domain

// Scorer represents a player in the historical scorers table.
type Scorer struct {
	PlayerID          int64      `json:"playerId"`
	FullName          string     `json:"fullName"`
	Team              SimpleTeam `json:"team"`
	Goals             int32      `json:"goals"`
	ListTeams         []string   `json:"listTeams"`
	ConfederationCode string     `json:"confederationCode"`
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
