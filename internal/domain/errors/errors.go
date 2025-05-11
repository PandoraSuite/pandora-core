package errors

type ErrorCode string

const (
	ErrorCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrorCodeInternal         ErrorCode = "INTERNAL"
	ErrorCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrorCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrorCodeAlreadyExists    ErrorCode = "ALREADY_EXISTS"
	ErrorCodeValidationFailed ErrorCode = "VALIDATION_FAILED"

	ErrorCodeAggregate ErrorCode = "AGGREGATE_ERRORS"
)

type Error interface {
	error

	Code() ErrorCode
}
