package errors

import "fmt"

var _ Error = (*BaseError)(nil)

type BaseError struct {
	code ErrorCode

	message string

	err error
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("<%s>: %s", e.code, e.message)
}

func (e *BaseError) Code() ErrorCode {
	return e.code
}

func (e *BaseError) Unwrap() error {
	return e.err
}

func NewNotFound(message string, err error) Error {
	return &BaseError{
		err:     err,
		code:    CodeNotFound,
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

func NewForbidden(message string, err error) Error {
	return &BaseError{
		err:     err,
		code:    CodeForbidden,
		message: message,
	}
}

func NewUnauthorized(message string, err error) Error {
	return &BaseError{
		err:     err,
		code:    CodeUnauthorized,
		message: message,
	}
}

func NewAlreadyExists(message string, err error) Error {
	return &BaseError{
		err:     err,
		code:    CodeAlreadyExists,
		message: message,
	}
}

func NewValidationFailed(message string, err error) Error {
	return &BaseError{
		err:     err,
		code:    CodeValidationFailed,
		message: message,
	}
}
