package errors

var (
	ErrRequestLogInvalidExecutionStatus = NewError(CodeValidationError, "Invalid Request Log execution status")
)
