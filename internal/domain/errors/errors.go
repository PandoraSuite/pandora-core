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

var ErrorCodePriority = map[ErrorCode]int{
	CodeInternal:         0,
	CodeUnauthorized:     1,
	CodeForbidden:        2,
	CodeNotFound:         3,
	CodeAlreadyExists:    4,
	CodeValidationFailed: 5,
	CodeAggregate:        6,
}

type Error interface {
	error

	Code() ErrorCode
	Unwrap() error
}
