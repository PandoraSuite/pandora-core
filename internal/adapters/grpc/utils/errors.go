package utils

import (
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
	"google.golang.org/grpc/codes"
)

func GetDomainErrorStatusCode(err *domainErr.Error) codes.Code {
	switch err.Code {
	case domainErr.CodeNotFound:
		return codes.NotFound
	case domainErr.CodeUnauthorized:
		return codes.PermissionDenied
	case domainErr.CodeInternalError:
		return codes.Internal
	case domainErr.CodeValidationError:
		return codes.FailedPrecondition
	default:
		return codes.Unknown
	}
}
