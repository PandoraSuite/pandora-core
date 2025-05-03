package errors

import "fmt"

var _ Error = (*BaseError)(nil)

type BaseError struct {
	code ErrorCode

	message string

	err error
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("<%s> %s", e.code, e.message)
}

func (e *BaseError) Code() ErrorCode {
	return e.code
}

func (e *BaseError) Unwrap() error {
	return e.err
}

func NewNotFound(message string) Error {
	return &BaseError{
		code:    CodeNotFound,
		message: message,
	}
}

func NewValidationFailed(message string) Error {
	return &BaseError{
		code:    CodeValidationFailed,
		message: message,
	}
}

func NewUnauthorized(message string) Error {
	return &BaseError{
		code:    CodeUnauthorized,
		message: message,
	}
}

func NewInternal(message string, err error) Error {
	return &BaseError{
		err:     err,
		code:    CodeInternal,
		message: message,
	}
}
