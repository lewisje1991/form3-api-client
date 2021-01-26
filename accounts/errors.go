package accounts

import "fmt"

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("status code %d: error %v", e.StatusCode, e.Message)
}
