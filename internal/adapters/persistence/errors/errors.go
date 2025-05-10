package errors

import (
	"fmt"
)

type ErrorCode string

const (
	ErrorCodeUnknown          ErrorCode = "UNKNOWN"
	ErrorCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrorCodeConnection       ErrorCode = "CONNECTION"
	ErrorCodeInvalidValue     ErrorCode = "INVALID_VALUE"
	ErrorCodeUniqueViolation  ErrorCode = "UNIQUE_VIOLATION"
	ErrorCodeUndefinedEntity  ErrorCode = "UNDEFINED_ENTITY"
	ErrorCodeInvalidReference ErrorCode = "INVALID_REFERENCE"
)

type Error interface {
	error

	Code() ErrorCode
	Field() string
	Entity() string
	Unwrap() error
}

type persistenceError struct {
	code ErrorCode

	field   string
	entity  string
	message string

	err error
}

func (e *persistenceError) Error() string {
	var entity string
	if e.entity != "" {
		entity = fmt.Sprintf("\n\tEntity: %s", e.entity)
	}

	var field string
	if e.field != "" {
		field = fmt.Sprintf("\n\tField: %s", e.field)
	}

	return fmt.Sprintf(
		"<%s>: %s%s\n\tMessage: %s",
		e.code,
		entity,
		field,
		e.message,
	)
}

func (e *persistenceError) Code() ErrorCode {
	return e.code
}

func (e *persistenceError) Field() string {
	return e.field
}

func (e *persistenceError) Entity() string {
	return e.entity
}

func (e *persistenceError) Unwrap() error {
	return e.err
}

func NewError(
	code ErrorCode, entity, field, message string, err error,
) *persistenceError {
	return &persistenceError{
		code:    code,
		field:   field,
		entity:  entity,
		message: message,
		err:     err,
	}
}

func NewNotFoundError(entity string, err error) *persistenceError {
	return &persistenceError{
		code:    ErrorCodeNotFound,
		entity:  entity,
		message: "Entity not found",
		err:     err,
	}
}

func NewConnectionError(message string, err error) *persistenceError {
	return &persistenceError{
		code:    ErrorCodeConnection,
		message: message,
		err:     err,
	}
}

func NewUnknownError(message string, err error) *persistenceError {
	return &persistenceError{
		code:    ErrorCodeUnknown,
		message: message,
		err:     err,
	}
}
