package validator

import (
	"fmt"
	"reflect"
	"strings"

	govalidator "github.com/go-playground/validator/v10"
)

const enumsTag = "enums"

func validateEnums(fl govalidator.FieldLevel) bool {
	allowedValuesParam := fl.Param()
	if allowedValuesParam == "" {
		return false
	}
	allowedValues := strings.Split(allowedValuesParam, " ")

	field := fl.Field()

	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			for _, allowed := range allowedValues {
				trimmedAllowed := strings.TrimSpace(allowed)
				if trimmedAllowed == "nil" || trimmedAllowed == "" {
					return true
				}
			}
			return false
		}
		field = field.Elem()
	}

	if !field.IsValid() {
		return false
	}

	fieldValueStr := fmt.Sprintf("%v", field.Interface())

	for _, allowed := range allowedValues {
		if fieldValueStr == strings.TrimSpace(allowed) {
			return true
		}
	}

	return false
}
