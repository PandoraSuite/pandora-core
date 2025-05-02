package errors

import "fmt"

var _ Error = (*EntityValidationError)(nil)

type EntityValidationError struct {
	BaseError

	entity string

	loc string
}

func (e *EntityValidationError) Error() string {
	var locPart string
	if e.loc != "" {
		locPart = fmt.Sprintf(" [%s]", e.loc)
	}

	return fmt.Sprintf(
		"<%s> %s%s: %s",
		e.code,
		e.entity,
		locPart,
		e.shortMsg,
	)
}

func NewEntityValidationError(entity, loc, msg string, err error) Error {
	return &EntityValidationError{
		BaseError: BaseError{
			code:     CodeValidationFailed,
			shortMsg: msg,
			err:      err,
		},
		entity: entity,
		loc:    loc,
	}
}
