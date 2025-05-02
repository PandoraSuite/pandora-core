package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type StructValidator struct {
	validator *validator.Validate
}

func (sv *StructValidator) ValidateStruct(value any, messages map[string]string) errors.Error {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return errors.NewInternal("cannot validate nil pointer", nil)
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.NewInternal(
			fmt.Sprintf(
				"input must be a non-nil struct or pointer to struct, got %T",
				value,
			),
			nil,
		)
	}

	var agg errors.Error

	errs := sv.validator.Struct(value)
	if errs != nil {
		validationErrors, ok := errs.(validator.ValidationErrors)
		if !ok {
			return errors.NewInternal(
				"unexpected error type from validator library", errs,
			)
		}

		for _, fieldErr := range validationErrors {
			loc := fieldErr.Namespace()

			structName, fieldPath := sv.separateStructNameAndPath(loc, val.Type())

			messageKey := fmt.Sprintf("%s.%s", fieldPath, fieldErr.Tag())
			message := sv.getMessage(messageKey, messages)

			agg = errors.Aggregate(
				agg,
				errors.NewEntityValidationError(
					structName, fieldPath, message, fieldErr,
				),
			)
		}
	}

	return agg
}

func (sv *StructValidator) getMessage(key string, messages map[string]string) string {
	re := regexp.MustCompile(`\[\d+\]`)
	normalizedKey := re.ReplaceAllString(key, "[]")

	if msg, ok := messages[normalizedKey]; ok {
		return msg
	}

	return fmt.Sprintf("Validation failed for field '%s'", key)
}

func (sv *StructValidator) separateStructNameAndPath(namespace string, rootType reflect.Type) (string, string) {
	structName := rootType.Name()
	if strings.HasPrefix(namespace, structName+".") {
		return structName, namespace[len(structName)+1:]
	}

	if namespace == structName {
		return structName, ""
	}

	return structName, namespace
}

func NewJSONStructValidator() *StructValidator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if jsonTag := fld.Tag.Get("json"); jsonTag != "" {
			return strings.Split(jsonTag, ",")[0]
		}
		return fld.Name
	})

	err := v.RegisterValidation(utcTag, validateUTC)
	if err != nil {
		panic(err)
	}

	return &StructValidator{validator: v}
}
