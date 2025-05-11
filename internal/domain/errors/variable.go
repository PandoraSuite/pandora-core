package errors

import "fmt"

var _ Error = (*VariableError)(nil)

type VariableError struct {
	BaseError

	name string
}

func (e *VariableError) Error() string {
	return fmt.Sprintf(
		"<%s> %s: %s",
		e.code,
		e.name,
		e.message,
	)
}

func NewVariableValidationFailed(name, message string, err error) Error {
	return &VariableError{
		BaseError: BaseError{
			err:     err,
			code:    ErrorCodeValidationFailed,
			message: message,
		},
		name: name,
	}
}
