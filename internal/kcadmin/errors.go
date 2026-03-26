package kcadmin

import "fmt"

// APIError represents a non-2xx Keycloak Admin API response.
type APIError struct {
	StatusCode int
	Path       string
	Body       string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("keycloak admin api error: status=%d path=%s body=%s", e.StatusCode, e.Path, e.Body)
}
