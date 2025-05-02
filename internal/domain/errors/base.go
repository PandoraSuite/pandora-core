package errors

import "fmt"

var _ Error = (*BaseError)(nil)

type BaseError struct {
	code ErrorCode

	shortMsg string

	err error
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("<%s> %s", e.code, e.shortMsg)
}

func (e *BaseError) Code() ErrorCode {
	return e.code
}

func (e *BaseError) Unwrap() error {
	return e.err
}
