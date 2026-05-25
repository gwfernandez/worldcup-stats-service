package domain

// NationalTeam represents a national team entity.
type NationalTeam struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	Code              string  `json:"code"`
	DissolutionDate   *string `json:"dissolution_date"`
	IsDissolved       bool    `json:"is_dissolved"`
	ConfederationCode string  `json:"confederation_code"`
	FederationName    string  `json:"federation_name"`
	FederationCode    string  `json:"federation_code"`
}

// NationalTeamFilter defines supported filters for listing national teams.
type NationalTeamFilter struct {
	Name              string
	ConfederationCode *string
	FederationName    string
	FederationCode    string
	IncludeDissolved  bool
	Page              int
	Size              int
}

// NationalTeamListResponse represents paginated national teams response.
type NationalTeamListResponse struct {
	Data       []NationalTeam `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}
