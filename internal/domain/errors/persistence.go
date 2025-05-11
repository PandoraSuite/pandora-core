package errors

import "fmt"

type PersistenceErrorCode string

const (
	PersistenceErrorCodeUnknown          PersistenceErrorCode = "UNKNOWN"
	PersistenceErrorCodeNotFound         PersistenceErrorCode = "NOT_FOUND"
	PersistenceErrorCodeConnection       PersistenceErrorCode = "CONNECTION"
	PersistenceErrorCodeInvalidValue     PersistenceErrorCode = "INVALID_VALUE"
	PersistenceErrorCodeUniqueViolation  PersistenceErrorCode = "UNIQUE_VIOLATION"
	PersistenceErrorCodeUndefinedEntity  PersistenceErrorCode = "UNDEFINED_ENTITY"
	PersistenceErrorCodeInvalidReference PersistenceErrorCode = "INVALID_REFERENCE"
)

type PersistenceError interface {
	error

	Code() PersistenceErrorCode
	Field() string
	Entity() string
	Unwrap() error
}

type persistenceError struct {
	code PersistenceErrorCode

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

func (e *persistenceError) Code() PersistenceErrorCode {
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

func NewPersistenceError(
	code PersistenceErrorCode, entity, field, message string, err error,
) *persistenceError {
	return &persistenceError{
		code:    code,
		field:   field,
		entity:  entity,
		message: message,
		err:     err,
	}
}

func NewPersistenceNotFoundError(entity string, err error) *persistenceError {
	return &persistenceError{
		code:    PersistenceErrorCodeNotFound,
		entity:  entity,
		message: "Entity not found",
		err:     err,
	}
}

func NewPersistenceConnectionError(message string, err error) *persistenceError {
	return &persistenceError{
		code:    PersistenceErrorCodeConnection,
		message: message,
		err:     err,
	}
}

func NewPersistenceUnknownError(message string, err error) *persistenceError {
	return &persistenceError{
		code:    PersistenceErrorCodeUnknown,
		message: message,
		err:     err,
	}
}
