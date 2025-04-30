package errors

type ErrorCode string

const (
	CodeNotFound         ErrorCode = "NOT_FOUND"
	CodeInternal         ErrorCode = "INTERNAL"
	CodeForbidden        ErrorCode = "FORBIDDEN"
	CodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	CodeAlreadyExists    ErrorCode = "ALREADY_EXISTS"
	CodeValidationFailed ErrorCode = "VALIDATION_FAILED"

	codeAggregate ErrorCode = "AGGREGATE"
)

type SimpleError struct {
	code ErrorCode

	shortMsg string

	err error
}

func (e *SimpleError) Error() string {
	return e.shortMsg
}

func (e *SimpleError) Unwrap() error {
	return e.err
}

func (e *SimpleError) Code() ErrorCode {
	return e.code
}
