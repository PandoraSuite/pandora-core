package errors

import (
	"net/http"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

func CodeToStatusCode(code errors.ErrorCode) int {
	switch code {
	case errors.CodeNotFound:
		return http.StatusNotFound
	case errors.CodeInternal:
		return http.StatusInternalServerError
	case errors.CodeForbidden:
		return http.StatusForbidden
	case errors.CodeUnauthorized:
		return http.StatusUnauthorized
	case errors.CodeAlreadyExists:
		return http.StatusConflict
	case errors.CodeValidationFailed:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
