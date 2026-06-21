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

// ScorerDetail represents a scorer with personal data and all valid goals.
type ScorerDetail struct {
	ID            int64        `json:"id"`
	FirstName     string       `json:"firstName"`
	LastName      string       `json:"lastName"`
	Position      *string      `json:"position"`
	Championships []int32      `json:"championships"`
	Teams         []SimpleTeam `json:"teams"`
	Goals         []Goal       `json:"goals"`
}
