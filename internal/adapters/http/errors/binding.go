package errors

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func BindingToHTTPError(structType any, err error) *HTTPError {
	switch e := err.(type) {
	case validator.ValidationErrors:
		errs := make([]*HTTPError, len(e))
		for i, ve := range e {
			errs[i] = NewValidationFailed(
				"body", getFieldName(structType, ve.Field()), ve.Error(),
			)
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

func getFieldName(strcutType any, fieldName string) string {
	t := reflect.TypeOf(strcutType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, ok := t.FieldByName(fieldName)
	if !ok {
		return fieldName
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName
	}

	jsonTagParts := strings.Split(jsonTag, ",")
	if len(jsonTagParts) == 0 {
		return fieldName
	}

	jsonFieldName := jsonTagParts[0]
	if jsonFieldName == "" {
		return fieldName
	}

	if jsonFieldName == "-" {
		return fieldName
	}

	return jsonFieldName
}
