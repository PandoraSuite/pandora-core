package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ErrorResponse struct {
	Error error `json:"error"`
}

func GetBindJSONErrorStatusCode(err error) int {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	if errors.As(err, &syntaxError) || errors.As(err, &unmarshalTypeError) {
		return http.StatusBadRequest
	}

	return http.StatusUnprocessableEntity
}

func GetDomainErrorStatusCode(err domainErr.Error) int {
	switch err.Code() {
	case domainErr.CodeNotFound:
		return http.StatusNotFound
	case domainErr.CodeForbidden:
		return http.StatusForbidden
	case domainErr.CodeUnauthorized:
		return http.StatusUnauthorized
	case domainErr.CodeInternal:
		return http.StatusInternalServerError
	case domainErr.CodeValidationFailed:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
