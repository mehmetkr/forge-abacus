package handler

// CalculateRequest represents the JSON request body.
type CalculateRequest struct {
	Operation string   `json:"operation"`
	A         *float64 `json:"a"`
	B         *float64 `json:"b"`
}

// CalculateResponse represents a successful computation result.
type CalculateResponse struct {
	Result float64 `json:"result"`
}

// ErrorResponse represents an error returned to the client.
type ErrorResponse struct {
	Error string `json:"error"`
}
