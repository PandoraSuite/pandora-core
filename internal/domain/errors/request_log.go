package errors

var (
	ErrRequestLogNotFound = NewError(CodeNotFound, "Request Log not found")

	ErrRequestLogInvalidExecutionStatus          = NewError(CodeValidationError, "Invalid Request Log execution status")
	ErrCannotUpdateToNullExecutionStatus         = NewError(CodeValidationError, "cannot update request log to null status")
	ErrCannotUpdateToPendingExecutionStatus      = NewError(CodeValidationError, "cannot update request log to pending status")
	ErrCannotUpdateToUnauthorizedExecutionStatus = NewError(CodeValidationError, "cannot update request log to unauthorized status")
)
