package validator

import (
	"fmt"
	"reflect"
	"time"

	govalidator "github.com/go-playground/validator/v10"
)

const utcTag = "utc"

func validateUTC(fl govalidator.FieldLevel) bool {
	field := fl.Field()

	var t time.Time
	isTime := false

	switch field.Kind() {
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			t = field.Interface().(time.Time)
			isTime = true
		}
	case reflect.Ptr:
		if field.Type().Elem() == reflect.TypeOf(time.Time{}) {
			if field.IsNil() {
				return true
			}
			t = field.Elem().Interface().(time.Time)
			isTime = true
		}
	}

	if !isTime {
		panic(
			fmt.Sprintf(
				"invalid usage of '%s' validation tag on field of type %s",
				utcTag, field.Type(),
			),
		)
	}

	return t.IsZero() || t.Location() == time.UTC
}
