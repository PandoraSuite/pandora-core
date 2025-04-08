package errors

var (
	ErrServiceNotFound             = NewError(CodeNotFound, "Service not found")
	ErrServiceNotAssignedToProject = NewError(CodeNotFound, "Service not assigned to project")

	ErrServiceDeprecated                   = NewError(CodeValidationError, "Service is deprecated")
	ErrServiceDeactivated                  = NewError(CodeValidationError, "Service is deactivated")
	ErrServiceInvalidStatus                = NewError(CodeValidationError, "Invalid Service status")
	ErrServiceAlreadyExistsWhitNameVersion = NewError(CodeValidationError, "Service already exists with that name and version")
)
