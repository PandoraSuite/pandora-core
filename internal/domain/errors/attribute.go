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

func (e *AttributeError) Entity() string {
	return e.entity
}

func (e *AttributeError) Loc() string {
	return e.loc
}

func (e *AttributeError) PrefixLoc(prefix string) {
	if e.loc != "" {
		e.loc = fmt.Sprintf("%s.%s", prefix, e.loc)
	} else {
		e.loc = prefix
	}
}

func NewAttributeNotFound(entity, loc, message string, err error) Error {
	return &AttributeError{
		BaseError: BaseError{
			err:     err,
			code:    CodeNotFound,
			message: message,
		},
		loc:    loc,
		entity: entity,
	}
}

func NewAttributeAlreadyExists(entity, loc, message string, err error) Error {
	return &AttributeError{
		BaseError: BaseError{
			err:     err,
			code:    CodeAlreadyExists,
			message: message,
		},
		loc:    loc,
		entity: entity,
	}
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
