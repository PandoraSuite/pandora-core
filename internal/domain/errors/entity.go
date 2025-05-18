package errors

import (
	"fmt"
	"sort"
	"strings"
)

var _ Error = (*EntityError)(nil)

type EntityError struct {
	BaseError

	entity string

	identifiers map[string]any
}

func (e *EntityError) Error() string {
	var idPart string
	if len(e.identifiers) > 0 {
		keys := make([]string, 0, len(e.identifiers))
		for k := range e.identifiers {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var b strings.Builder
		for i, k := range keys {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "%s=%v", k, e.identifiers[k])
		}
		idPart = fmt.Sprintf(" (%s)", b.String())
	}

	return fmt.Sprintf(
		"<%s> %s%s: %s",
		e.code,
		e.entity,
		idPart,
		e.message,
	)
}

func (e *EntityError) Entity() string {
	return e.entity
}

func (e *EntityError) Identifiers() map[string]any {
	return e.identifiers
}

func NewEntityNotFound(
	entity, message string, identifiers map[string]any, err error,
) *EntityError {
	return &EntityError{
		BaseError: BaseError{
			err:     err,
			code:    CodeNotFound,
			message: message,
		},
		entity:      entity,
		identifiers: identifiers,
	}
}

func NewEntityAlreadyExists(
	entity, message string, identifiers map[string]any, err error,
) *EntityError {
	return &EntityError{
		BaseError: BaseError{
			err:     err,
			code:    CodeAlreadyExists,
			message: message,
		},
		entity:      entity,
		identifiers: identifiers,
	}
}

func NewEntityValidationFailed(
	entity, message string, identifiers map[string]any, err error,
) *EntityError {
	return &EntityError{
		BaseError: BaseError{
			err:     err,
			code:    CodeValidationFailed,
			message: message,
		},
		entity:      entity,
		identifiers: identifiers,
	}
}
