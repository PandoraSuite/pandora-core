package errors

import (
	"google.golang.org/grpc/codes"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

func CodeToGRPCCode(code errors.ErrorCode) codes.Code {
	switch code {
	case errors.CodeNotFound:
		return codes.NotFound
	case errors.CodeInternal:
		return codes.Internal
	case errors.CodeForbidden:
		return codes.PermissionDenied
	case errors.CodeUnauthorized:
		return codes.Unauthenticated
	case errors.CodeAlreadyExists:
		return codes.AlreadyExists
	case errors.CodeValidationFailed:
		return codes.InvalidArgument
	default:
		return codes.Unknown
	}
}
