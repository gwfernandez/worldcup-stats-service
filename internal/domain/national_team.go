package domain

// NationalTeam represents a national team entity.
type NationalTeam struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Code            string  `json:"code"`
	DissolutionDate *string `json:"dissolution_date"`
	IsDissolved     bool    `json:"is_dissolved"`
	ConfederationID int64   `json:"confederation_id"`
	FederationName  string  `json:"federation_name"`
	FederationCode  string  `json:"federation_code"`
}

// NationalTeamFilter defines supported filters for listing national teams.
type NationalTeamFilter struct {
	Name             string
	ConfederationID  *int64
	FederationName   string
	FederationCode   string
	IncludeDissolved bool
	Page             int
	Size             int
}

// PaginatedNationalTeams represents paginated national teams response.
type PaginatedNationalTeams struct {
	Content       []NationalTeam `json:"content"`
	Page          int            `json:"page"`
	Size          int            `json:"size"`
	TotalElements int64          `json:"total_elements"`
	TotalPages    int            `json:"total_pages"`
	HasNext       bool           `json:"has_next"`
	HasPrevious   bool           `json:"has_previous"`
}
