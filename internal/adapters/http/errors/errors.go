package errors

import "github.com/MAD-py/pandora-core/internal/domain/errors"

type HTTPError struct {
	Code    errors.ErrorCode `json:"code"`
	Message string           `json:"message"`

	Entity      string         `json:"entity,omitempty"`
	Identifiers map[string]any `json:"identifiers,omitempty"`

	Loc string `json:"loc,omitempty"`

	Errors []*HTTPError `json:"errors,omitempty"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) PriorityCode() errors.ErrorCode {
	if len(e.Errors) == 0 {
		return e.Code
	}

	best := e.Code
	bestPriority := errors.ErrorCodePriority[best]

	for _, err := range e.Errors {
		priority, ok := errors.ErrorCodePriority[err.Code]
		if ok && priority < bestPriority {
			best = err.Code
			bestPriority = priority
		}
	}

	return best
}

func NewValidationFailed(entity, loc, message string) *HTTPError {
	return &HTTPError{
		Code:    errors.CodeValidationFailed,
		Message: message,

		Entity: entity,
		Loc:    loc,
	}
}

func NewUnauthorized(message string) *HTTPError {
	return &HTTPError{
		Code:    errors.CodeUnauthorized,
		Message: message,
	}
}

func NewForbidden(message string) *HTTPError {
	return &HTTPError{
		Code:    errors.CodeForbidden,
		Message: message,
	}
}

func NewInternal(message string) *HTTPError {
	return &HTTPError{
		Code:    errors.CodeInternal,
		Message: message,
	}
}
