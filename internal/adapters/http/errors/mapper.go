package errors

import "github.com/MAD-py/pandora-core/internal/domain/errors"

func MapToHTTPError(err error) *HTTPError {
	switch e := err.(type) {
	case *HTTPError:
		return e
	case *errors.BaseError:
		return &HTTPError{
			Code:    e.Code(),
			Message: e.Message(),
		}
	case *errors.VariableError:
		return &HTTPError{
			Code:    e.Code(),
			Message: e.Message(),

			Loc: e.Name(),
		}
	case *errors.EntityError:
		return &HTTPError{
			Code:    e.Code(),
			Message: e.Message(),

			Entity:      e.Entity(),
			Identifiers: e.Identifiers(),
		}
	case *errors.AttributeError:
		return &HTTPError{
			Code:    e.Code(),
			Message: e.Message(),

			Entity: e.Entity(),
			Loc:    e.Loc(),
		}
	case errors.AggregateError:
		errs := make([]*HTTPError, len(e))
		for i, err := range e {
			errs[i] = MapToHTTPError(err)
		}

		return &HTTPError{
			Code:    e.Code(),
			Message: "Multiple errors occurred",

			Errors: errs,
		}
	default:
		return &HTTPError{
			Code:    errors.CodeInternal,
			Message: "Internal server error",
		}
	}
}
