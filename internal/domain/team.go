package domain

// Team represents a team entity.
type Team struct {
	Code              string  `json:"code"`
	Name              string  `json:"name"`
	IsDissolved       bool    `json:"isDissolved"`
	ConfederationCode string  `json:"confederationCode"`
	FederationName    string  `json:"federationName"`
	FederationCode    string  `json:"federationCode"`
	DissolutionDate   *string `json:"dissolutionDate"`
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
