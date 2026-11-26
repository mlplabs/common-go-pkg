package custom

import "net/http"

type ConflictError struct {
	code string
	err  error
}

func (e *ConflictError) StatusCode() int { return http.StatusConflict }
func (e *ConflictError) ErrorCode() string {
	if e.code == "" {
		return "CONFLICT"
	}
	return e.code
}
func (e *ConflictError) Error() string { return e.err.Error() }

func NewConflict(code string, err error) error {
	return &ConflictError{code: code, err: err}
}
