package validator

import (
	"fmt"
	"reflect"
	"time"

	govalidator "github.com/go-playground/validator/v10"
)

const utcTag = "utc"

func validateUTC(fl govalidator.FieldLevel) bool {
	t := convertToTime(fl.Field())
	return t.IsZero() || t.Location() == time.UTC
}

const gteTimeTag = "gtetime"

func validateGteTime(fl govalidator.FieldLevel) bool {
	t := convertToTime(fl.Field())

	if t.IsZero() {
		return true
	}

	otherFieldName := fl.Param()
	if otherFieldName == "" {
		panic(
			fmt.Sprintf(
				"invalid usage of '%s' validation tag without a field name",
				gteTimeTag,
			),
		)
	}

	otherT := convertToTime(fl.Parent().FieldByName(otherFieldName))

	if otherT.IsZero() {
		return false
	}

	return t.After(otherT) || t.Equal(otherT)
}

func convertToTime(field reflect.Value) time.Time {
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
	return t
}
