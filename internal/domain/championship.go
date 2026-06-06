package domain

// TopScorer represents a top scorer of a championship edition.
type TopScorer struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	NationCode string `json:"nation_code"`
}

// ChampionshipsStats represents statistics for a specific championship edition.
type ChampionshipsStats struct {
	TotalTeams      int32       `json:"total_teams"`
	TotalMatches    int32       `json:"total_matches"`
	TotalStadiums   int32       `json:"total_stadiums"`
	TotalPlayers    int32       `json:"total_players"`
	TotalGoals      int32       `json:"total_goals"`
	RunnerUpCode    string      `json:"runner_up_code"`
	ThirdPlaceCode  string      `json:"third_place_code"`
	FourthPlaceCode string      `json:"fourth_place_code"`
	TopScorers      []TopScorer `json:"top_scorers"`
	TopScorerGoals  int32       `json:"top_scorer_goals"`
}

// Championship represents a football world cup championship edition.
type Championship struct {
	Year         int                 `json:"year"`
	StartDate    string              `json:"start_date"`
	EndDate      string              `json:"end_date"`
	HostCodes    []string            `json:"host_codes"`
	ChampionCode *string             `json:"champion_code"`
	Stats        *ChampionshipsStats `json:"stats,omitempty"`
}

// ChampionshipFilter represents filters for listing championships.
type ChampionshipFilter struct {
	Year              int
	Host              string
	ConfederationCode string
	Page              int
	Size              int
}

// ChampionshipTeam represents a team that participated in a championship edition.
type ChampionshipTeam struct {
	Year              int    `json:"year"`
	TeamCode          string `json:"team_code"`
	ConfederationCode string `json:"confederation_code"`
	GroupCode         string `json:"group_code"`
	StageReached      string `json:"stage_reached"`
	Managers          string `json:"managers"`
}

// ChampionshipTeamFilter represents filters for listing championship teams.
type ChampionshipTeamFilter struct {
	Year              int
	Name              string
	ConfederationCode string
	GroupCode         string
	Page              int
	Size              int
}

// ChampionshipListResponse represents the JSON response for listing championships.
type ChampionshipListResponse struct {
	Data       []Championship `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

// ChampionshipTeamListResponse represents the JSON response for listing championship teams.
type ChampionshipTeamListResponse struct {
	Data       []ChampionshipTeam `json:"data"`
	Pagination PaginationInfo     `json:"pagination"`
}
