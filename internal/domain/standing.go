package domain

// Standing represents a team in the historical standings table.
type Standing struct {
	Team            SimpleTeam `json:"team"`
	MatchesPlayed   int32      `json:"matchesPlayed"`
	Wins            int32      `json:"wins"`
	Draws           int32      `json:"draws"`
	Losses          int32      `json:"losses"`
	GoalsFor        int32      `json:"goalsFor"`
	GoalsAgainst    int32      `json:"goalsAgainst"`
	GoalDifference  int32      `json:"goalDifference"`
	Points          int32      `json:"points"`
	UnifiedPoints   int32      `json:"unifiedPoints"`
	Position        int32      `json:"position"`
	UnifiedPosition int32      `json:"unifiedPosition"`
}

// StandingFilter defines filters and pagination for listing historical standings.
type StandingFilter struct {
	Name              string
	Language          string
	ConfederationCode string
	Page              int
	Size              int
}

// StandingListResponse represents paginated historical standings response.
type StandingListResponse struct {
	Data       []Standing     `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}
