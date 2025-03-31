package errors

import "fmt"

type ErrorCode string

const (
	CodeNotFound        ErrorCode = "NOT_FOUND"
	CodeForbidden       ErrorCode = "FORBIDDEN"
	CodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	CodeInternalError   ErrorCode = "INTERNAL_ERROR"
	CodeValidationError ErrorCode = "VALIDATION_ERROR"
)

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details []string  `json:"details,omitempty"`
}

func (e *Error) AddDetail(detail string) *Error {
	copyErr := *e

	copyErr.Details = make([]string, len(e.Details))
	copy(copyErr.Details, e.Details)

	copyErr.Details = append(copyErr.Details, detail)

	return &copyErr

}

func (e *Error) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("(%s) %s: %v", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("(%s) %s", e.Code, e.Message)
}

func NewError(code ErrorCode, message string, details ...string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}
