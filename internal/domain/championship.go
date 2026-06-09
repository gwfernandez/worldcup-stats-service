package domain

// TopScorer represents a top scorer of a championship edition.
type TopScorer struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	NationCode string `json:"nationCode"`
}

// ChampionshipsStats represents statistics for a specific championship edition.
type ChampionshipsStats struct {
	TotalTeams      int32       `json:"totalTeams"`
	TotalMatches    int32       `json:"totalMatches"`
	TotalStadiums   int32       `json:"totalStadiums"`
	TotalPlayers    int32       `json:"totalPlayers"`
	TotalGoals      int32       `json:"totalGoals"`
	RunnerUpCode    string      `json:"runnerUpCode"`
	ThirdPlaceCode  string      `json:"thirdPlaceCode"`
	FourthPlaceCode string      `json:"fourthPlaceCode"`
	TopScorers      []TopScorer `json:"topScorers"`
	TopScorerGoals  int32       `json:"topScorerGoals"`
}

// Championship represents a football world cup championship edition.
type Championship struct {
	Year         int                 `json:"year"`
	StartDate    string              `json:"startDate"`
	EndDate      string              `json:"endDate"`
	HostCodes    []string            `json:"hostCodes"`
	ChampionCode *string             `json:"championCode"`
	ChampionName *string             `json:"championName,omitempty"`
	Stats        *ChampionshipsStats `json:"stats,omitempty"`
}

// ChampionshipFilter represents filters for listing championships.
type ChampionshipFilter struct {
	Year              int
	Host              string
	Language          string
	ConfederationCode string
	Page              int
	Size              int
}

// ChampionshipTeam represents a team that participated in a championship edition.
type ChampionshipTeam struct {
	Year              int    `json:"year"`
	TeamCode          string `json:"teamCode"`
	Name              string `json:"name"`
	ConfederationCode string `json:"confederationCode"`
	GroupCode         string `json:"groupCode"`
	StageReached      string `json:"stageReached"`
	Managers          string `json:"managers"`
}

// ChampionshipTeamFilter represents filters for listing championship teams.
type ChampionshipTeamFilter struct {
	Year              int
	Name              string
	Language          string
	ConfederationCode string
	GroupCode         string
	Page              int
	Size              int
}

// ChampionshipStadium represents a stadium used in a championship edition.
type ChampionshipStadium struct {
	Year          int    `json:"year"`
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	CityName      string `json:"cityName"`
	Capacity      int32  `json:"capacity"`
	MatchesPlayed int32  `json:"matchesPlayed"`
}

// ChampionshipStadiumFilter represents filters for listing championship stadiums.
type ChampionshipStadiumFilter struct {
	Year int
	Name string
	Page int
	Size int
}

// ChampionshipScorer represents a player scorer in a championship edition.
type ChampionshipScorer struct {
	FullName string `json:"fullName"`
	TeamCode string `json:"teamCode"`
	TeamName string `json:"teamName"`
	Goals    int32  `json:"goals"`
}

// ChampionshipScorerFilter represents filters for listing championship scorers.
type ChampionshipScorerFilter struct {
	Year     int
	Name     string
	Language string
	TeamCode string
	Page     int
	Size     int
}

// ChampionshipStanding represents a team standing in a championship edition.
type ChampionshipStanding struct {
	TeamCode       string `json:"teamCode"`
	TeamName       string `json:"teamName"`
	GroupCode      string `json:"groupCode"`
	MatchesPlayed  int32  `json:"matchesPlayed"`
	Wins           int32  `json:"wins"`
	Draws          int32  `json:"draws"`
	Losses         int32  `json:"losses"`
	GoalsFor       int32  `json:"goalsFor"`
	GoalsAgainst   int32  `json:"goalsAgainst"`
	GoalDifference int32  `json:"goalDifference"`
	Points         int32  `json:"points"`
	UnifiedPoints  int32  `json:"unifiedPoints"`
	Position       int32  `json:"position"`
	Performance    string `json:"performance"`
}

// ChampionshipStandingFilter represents filters for listing championship standings.
type ChampionshipStandingFilter struct {
	Year     int
	Language string
	Page     int
	Size     int
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

// ChampionshipStadiumListResponse represents the JSON response for listing championship stadiums.
type ChampionshipStadiumListResponse struct {
	Data       []ChampionshipStadium `json:"data"`
	Pagination PaginationInfo        `json:"pagination"`
}

// ChampionshipScorerListResponse represents the JSON response for listing championship scorers.
type ChampionshipScorerListResponse struct {
	Data       []ChampionshipScorer `json:"data"`
	Pagination PaginationInfo       `json:"pagination"`
}

// ChampionshipStandingListResponse represents the JSON response for listing championship standings.
type ChampionshipStandingListResponse struct {
	Data       []ChampionshipStanding `json:"data"`
	Pagination PaginationInfo         `json:"pagination"`
}
