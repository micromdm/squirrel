package models

// ErrorResponse encodes errors into http response body
type ErrorResponse struct {
	Errors []string `json:"errors" plist:"errors"`
}

// View returns a view
func (e *ErrorResponse) View(accept string) (*Response, error) {
	return marshal(e, accept)
}
