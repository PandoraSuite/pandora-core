package utils

import (
	"net/http"

	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ErrorResponse struct {
	Error error `json:"error"`
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
