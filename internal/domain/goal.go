package domain

// Goal represents a goal scored by a player in a World Cup match.
type Goal struct {
	Year          int32        `json:"year"`
	Hosts         []SimpleTeam `json:"hosts"`
	MatchDate     *string      `json:"matchDate"`
	OpponentTeam  SimpleTeam   `json:"opponentTeam"`
	MinuteRegular int32        `json:"minuteRegular"`
	Penalty       *bool        `json:"penalty"`
	Stage         *string      `json:"stage"`
}

// GoalFilter defines player, championship and pagination filters for goals.
type GoalFilter struct {
	PlayerID int64
	Year     int
	Language string
	Page     int
	Size     int
}

// GoalListResponse represents a paginated player goals response.
type GoalListResponse struct {
	Data       []Goal         `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}
