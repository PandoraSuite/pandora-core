package errors

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func BindingToHTTPError(err error) *HTTPError {
	switch e := err.(type) {
	case validator.ValidationErrors:
		errs := make([]*HTTPError, len(e))
		for i, ve := range e {
			errs[i] = NewValidationFailed("body", ve.Field(), ve.Error())
		}

		if len(errs) == 1 {
			return errs[0]
		}

		return NewMultipleErrors(errs)

	case *json.UnmarshalTypeError:
		return NewValidationFailed(
			"body",
			e.Field,
			fmt.Sprintf(
				"Invalid type for field '%s', expected %s",
				e.Field, e.Type.String(),
			),
		)
	case *json.SyntaxError:
		return NewValidationFailed(
			"body",
			fmt.Sprint(e.Offset),
			fmt.Sprintf(
				"Malformed JSON in request body, offset: %d", e.Offset,
			),
		)
	default:
		return NewValidationFailed("body", "", "Invalid request payload")
	}
}
