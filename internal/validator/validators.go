package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	govalidator "github.com/go-playground/validator/v10"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type Validator interface {
	ValidateStruct(value any, messages map[string]string) errors.Error
	ValidateVariable(value any, fieldName, tags string, messages map[string]string) errors.Error
}

type validator struct {
	validator *govalidator.Validate
}

func (sv *validator) ValidateStruct(value any, messages map[string]string) errors.Error {
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
		validationErrors, ok := errs.(govalidator.ValidationErrors)
		if !ok {
			return errors.NewInternal(
				"unexpected error type from validator", errs,
			)
		}

		for _, fieldErr := range validationErrors {
			loc := fieldErr.Namespace()

			structName, fieldPath := sv.separateStructNameAndPath(loc, val.Type())

			messageKey := fmt.Sprintf("%s.%s", fieldPath, fieldErr.Tag())
			message := sv.getMessage(messageKey, messages)

			agg = errors.Aggregate(
				agg,
				errors.NewAttributeValidationFailed(
					structName, fieldPath, message, fieldErr,
				),
			)
		}
	}

	return agg
}

func (sv *validator) ValidateVariable(value any, fieldName, tags string, messages map[string]string) errors.Error {
	var agg errors.Error

	errs := sv.validator.Var(value, tags)

	if errs != nil {
		validationErrors, ok := errs.(govalidator.ValidationErrors)
		if !ok {
			return errors.NewInternal(
				"unexpected error type from validator", errs,
			)
		}

		for _, fieldErr := range validationErrors {
			failedTag := fieldErr.Tag()

			message := sv.getMessage(failedTag, messages)

			agg = errors.Aggregate(
				agg,
				errors.NewVariableValidationFailed(
					fieldName, message, fieldErr,
				),
			)
		}
	}

	return agg
}

func (sv *validator) getMessage(key string, messages map[string]string) string {
	re := regexp.MustCompile(`\[\d+\]`)
	normalizedKey := re.ReplaceAllString(key, "[]")

	if msg, ok := messages[normalizedKey]; ok {
		return msg
	}

	return fmt.Sprintf("Validation failed for field '%s'", key)
}

func (sv *validator) separateStructNameAndPath(namespace string, rootType reflect.Type) (string, string) {
	structName := rootType.Name()
	if strings.HasPrefix(namespace, structName+".") {
		return structName, namespace[len(structName)+1:]
	}

	if namespace == structName {
		return structName, ""
	}

	return structName, namespace
}

func NewValidator() Validator {
	v := govalidator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if name := fld.Tag.Get("name"); name != "" {
			return name
		}
		return strings.ToLower(fld.Name)
	})

	err := v.RegisterValidation(utcTag, validateUTC)
	if err != nil {
		panic(err)
	}

	err = v.RegisterValidation(enumsTag, validateEnums)
	if err != nil {
		panic(err)
	}

	return &validator{validator: v}
}
