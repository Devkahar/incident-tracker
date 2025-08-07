package errors

import "fmt"

type APIError interface {
	error
	StatusCode() int
	Json() map[string]any
}
type _APIError struct {
	Status  int   `json:"status"`
	Message any   `json:"message"`
	Cause   error `json:"error"`
}

func NewAPIError(statusCode int, message any, err error) APIError {
	return APIError(&_APIError{
		Status:  statusCode,
		Message: message,
		Cause:   err,
	})
}

func (_e *_APIError) Error() string {
	if _e.Cause != nil {
		return fmt.Sprintf("status=%d, message=%v, cause=%v", _e.Status, _e.Message, _e.Cause.Error())
	}
	return fmt.Sprintf("status=%d, message=%v", _e.Status, _e.Message)
}

func (e *_APIError) Json() map[string]any {
	return map[string]any{
		"status":  e.Status,
		"message": e.Message,
		"error":   e.Cause,
	}
}

func (e *_APIError) StatusCode() int {
	return e.Status
}
