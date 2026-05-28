package domain

// Team represents a team entity.
type Team struct {
	Code              string  `json:"code"`
	Name              string  `json:"name"`
	IsDissolved       bool    `json:"is_dissolved"`
	ConfederationCode string  `json:"confederation_code"`
	FederationName    string  `json:"federation_name"`
	FederationCode    string  `json:"federation_code"`
	DissolutionDate   *string `json:"dissolution_date"`
}

// TeamFilter defines supported filters for listing teams.
type TeamFilter struct {
	Name              string
	ConfederationCode *string
	FederationName    string
	FederationCode    string
	IncludeDissolved  bool
	Page              int
	Size              int
}

// TeamListResponse represents paginated teams response.
type TeamListResponse struct {
	Data       []Team         `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}
