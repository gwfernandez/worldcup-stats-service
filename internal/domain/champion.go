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

// ChampionFinal represents a World Cup final won by a team.
type ChampionFinal struct {
	Year                   int32      `json:"year"`
	MatchDate              *string    `json:"matchDate"`
	MatchTime              *string    `json:"matchTime"`
	HomeTeam               SimpleTeam `json:"homeTeam"`
	HomeTeamScore          *int32     `json:"homeTeamScore"`
	HomeTeamScorePenalties *int32     `json:"homeTeamScorePenalties"`
	AwayTeam               SimpleTeam `json:"awayTeam"`
	AwayTeamScore          *int32     `json:"awayTeamScore"`
	AwayTeamScorePenalties *int32     `json:"awayTeamScorePenalties"`
}

// ChampionFinalFilter defines pagination and team selection for won finals.
type ChampionFinalFilter struct {
	TeamCode string
	Language string
	Page     int
	Size     int
}

// ChampionFinalListResponse represents paginated finals won by a team.
type ChampionFinalListResponse struct {
	Data       []ChampionFinal `json:"data"`
	Pagination PaginationInfo  `json:"pagination"`
}
