package domain

// Confederation represents a football confederation entity.
type Confederation struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// CreateConfederationRequest holds the input data for creating a confederation.
type CreateConfederationRequest struct {
	Code string `json:"code" binding:"required,max=20"`
	Name string `json:"name" binding:"required,max=100"`
}

// UpdateConfederationRequest holds the input data for updating a confederation.
type UpdateConfederationRequest struct {
	Code string `json:"code" binding:"required,max=20"`
	Name string `json:"name" binding:"required,max=100"`
}
