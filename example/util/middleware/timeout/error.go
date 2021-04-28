package timeout

import (
	"errors"
	"fmt"
	"net/http"
)

type Type string

const (
	Authorization        Type = "AUTHORIZATION"
	BadRequest           Type = "BAD_REQUEST"
	Conflict             Type = "CONFLICT"
	Internal             Type = "INTERNAL"
	NotFound             Type = "NOT_FOUND"
	PayloadTooLarge      Type = "PAYLOAD_TOO_LARGE"
	ServiceUnavailable   Type = "SERVICE_UNAVAILABLE"
	UnsupportedMediaType Type = "UNSUPPORTED_MEDIA_TYPE"
	RequestTimeout       Type = "REQUEST_TIMEOUT"
)

// Error holds a custom error for the application
// which is helpful in returning a consistent
// error type/message from API endpoints
type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

// Error satisfies standard error interface
// we can return errors from this package as
// a regular old go _error_
func (e *Error) Error() string {
	return e.Message
}

// Status is a mapping errors to status codes
// Of course, this is somewhat redundant since
// our errors already map http status codes
func (e *Error) Status() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case RequestTimeout:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}

// Status checks the runtime type
// of the error and returns an http
// status code if the error is model.Error
func Status(err error) int {
	var e *Error

	if errors.As(err, e) {
		return e.Status()
	}

	return http.StatusInternalServerError
}

// NewAuthorization to create a 401
func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

// NewBadRequest to create 400 errors (validation, for example)
func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: reason,
	}
}

// NewConflict to create an error for 409
func NewConflict(name string, value string) *Error {
	return &Error{
		Type:    Conflict,
		Message: fmt.Sprintf("resource: %v with value: %v already exists", name, value),
	}
}

// NewInternal for 500 errors and unknown errors
func NewInternal() *Error {
	return &Error{
		Type:    Internal,
		Message: "Internal server error.",
	}
}

// NewNotFound to create an error for 404
func NewNotFound(name string, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("resource: %v with value: %v not found", name, value),
	}
}

// NewPayloadTooLarge to create an error for 413
func NewPayloadTooLarge(max int64, contentLength int64) *Error {
	return &Error{
		Type:    PayloadTooLarge,
		Message: fmt.Sprintf("Max payload size of %v exceeded. Actual payload size: %v", max, contentLength),
	}
}

// NewServiceUnavailable to create an error for 503
func NewServiceUnavailable() *Error {
	return &Error{
		Type:    ServiceUnavailable,
		Message: "Service unavailable or timeout",
	}
}

// NewUnsupportedMediaType to create an error for 415
func NewUnsupportedMediaType(reason string) *Error {
	return &Error{
		Type:    UnsupportedMediaType,
		Message: reason,
	}
}

// NewRequestTimeout to create an error for 408
func NewRequestTimeout(reason string) *Error {
	return &Error{
		Type:    RequestTimeout,
		Message: reason,
	}
}
