package errors

import "fmt"

var _ Error = (*AttributeError)(nil)

type AttributeError struct {
	BaseError

	entity string

	loc string
}

func (e *AttributeError) Error() string {
	var locPart string
	if e.loc != "" {
		locPart = fmt.Sprintf(" [%s]", e.loc)
	}

	return fmt.Sprintf(
		"<%s> %s%s: %s",
		e.code,
		e.entity,
		locPart,
		e.message,
	)
}

func NewAttributeValidationFailed(entity, loc, message string, err error) Error {
	return &AttributeError{
		BaseError: BaseError{
			err:     err,
			code:    CodeValidationFailed,
			message: message,
		},
		loc:    loc,
		entity: entity,
	}
}
