package errors

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

func BindingToDomainError(err error) errors.Error {
	switch e := err.(type) {
	case validator.ValidationErrors:
		var domainErr errors.Error
		for _, ve := range e {
			errors.Aggregate(
				domainErr,
				errors.NewAttributeValidationFailed(
					"body",
					ve.Field(),
					ve.Error(),
					ve,
				),
			)
		}
		return domainErr
	case *json.UnmarshalTypeError:
		return errors.NewAttributeValidationFailed(
			"body",
			e.Field,
			fmt.Sprintf(
				"Invalid type for field '%s', expected %s",
				e.Field, e.Type.String(),
			),
			e,
		)
	case *json.SyntaxError:
		return errors.NewValidationFailed(
			fmt.Sprintf(
				"Malformed JSON in request body, offset: %d", e.Offset,
			),
			e,
		)
	default:
		return errors.NewValidationFailed(
			"Invalid request payload.", err,
		)
	}
}
