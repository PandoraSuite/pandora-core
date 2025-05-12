package utils

import (
	"net/http"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ErrorResponse struct {
	Error error `json:"error"`
}

func GetDomainErrorStatusCode(err errors.Error) int {
	switch err.Code() {
	case errors.CodeNotFound:
		return http.StatusNotFound
	case errors.CodeForbidden:
		return http.StatusForbidden
	case errors.CodeUnauthorized:
		return http.StatusUnauthorized
	case errors.CodeInternal:
		return http.StatusInternalServerError
	case errors.CodeValidationFailed:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
