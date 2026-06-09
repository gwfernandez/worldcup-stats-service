package domain

// FixtureMatch represents a match in a championship fixture response.
type FixtureMatch struct {
	ID                     int64   `json:"id"`
	StageType              string  `json:"stageType"`
	Replayed               bool    `json:"replayed"`
	ReplayOf               *int64  `json:"replayOf"`
	MatchDate              *string `json:"matchDate"`
	MatchTime              *string `json:"matchTime"`
	StadiumID              *int64  `json:"stadiumId"`
	HomeTeamCode           string  `json:"homeTeamCode"`
	HomeTeamName           string  `json:"homeTeamName"`
	AwayTeamCode           string  `json:"awayTeamCode"`
	AwayTeamName           string  `json:"awayTeamName"`
	HomeTeamScore          *int32  `json:"homeTeamScore"`
	AwayTeamScore          *int32  `json:"awayTeamScore"`
	ExtraTime              bool    `json:"extraTime"`
	PenaltyShootout        bool    `json:"penaltyShootout"`
	HomeTeamScorePenalties *int32  `json:"homeTeamScorePenalties"`
	AwayTeamScorePenalties *int32  `json:"awayTeamScorePenalties"`
	HomeTeamWin            *bool   `json:"homeTeamWin"`
	AwayTeamWin            *bool   `json:"awayTeamWin"`
	Draw                   *bool   `json:"draw"`
	RefID                  *string `json:"refId"`
}

// GroupStanding represents a team standing inside a fixture group.
type GroupStanding struct {
	TeamCode       string `json:"teamCode"`
	TeamName       string `json:"teamName"`
	MatchesPlayed  int32  `json:"matchesPlayed"`
	Wins           int32  `json:"wins"`
	Draws          int32  `json:"draws"`
	Losses         int32  `json:"losses"`
	GoalsFor       int32  `json:"goalsFor"`
	GoalsAgainst   int32  `json:"goalsAgainst"`
	GoalDifference int32  `json:"goalDifference"`
	Points         int32  `json:"points"`
	UnifiedPoints  int32  `json:"unifiedPoints"`
	Position       *int32 `json:"position"`
}

// FixtureGroup represents a group stage bucket with matches and standings.
type FixtureGroup struct {
	GroupCode string          `json:"groupCode"`
	Matches   []FixtureMatch  `json:"matches"`
	Standings []GroupStanding `json:"standings"`
}

// FixtureStage represents a stage in a championship fixture.
type FixtureStage struct {
	Stage   string         `json:"stage"`
	Groups  []FixtureGroup `json:"groups,omitempty"`
	Matches []FixtureMatch `json:"matches,omitempty"`
}

// Fixture represents the full fixture for a championship year.
type Fixture struct {
	Year   int            `json:"year"`
	Stages []FixtureStage `json:"stages"`
}

// FixtureMatchRecord carries a fixture match with grouping metadata.
type FixtureMatchRecord struct {
	Stage     string
	GroupCode string
	Match     FixtureMatch
}

// GroupStandingRecord carries a group standing with grouping metadata.
type GroupStandingRecord struct {
	Stage     string
	GroupCode string
	Standing  GroupStanding
}
