package errors

type ErrorCode string

const (
	CodeNotFound         ErrorCode = "NOT_FOUND"
	CodeInternal         ErrorCode = "INTERNAL"
	CodeForbidden        ErrorCode = "FORBIDDEN"
	CodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	CodeAlreadyExists    ErrorCode = "ALREADY_EXISTS"
	CodeValidationFailed ErrorCode = "VALIDATION_FAILED"

	CodeAggregate ErrorCode = "AGGREGATE_ERRORS"
)

type Error interface {
	error

	Code() ErrorCode
}
