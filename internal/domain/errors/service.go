package errors

var (
	ErrServiceNotFound                 = NewError(CodeNotFound, "Service not found")
	ErrServiceNotAssignedToProject     = NewError(CodeNotFound, "Service not assigned to Project")
	ErrServiceNotAssignedToEnvironment = NewError(CodeNotFound, "Service not assigned to Environment")

	ErrServiceDeprecated                   = NewError(CodeValidationError, "Service is deprecated")
	ErrServiceDeactivated                  = NewError(CodeValidationError, "Service is deactivated")
	ErrServiceInvalidStatus                = NewError(CodeValidationError, "Invalid Service status")
	ErrServiceAssignedToProjects           = NewError(CodeValidationError, "Cannot delete Service because it is assigned to one or more Projects")
	ErrServiceAlreadyExistsWhitNameVersion = NewError(CodeValidationError, "Service already exists with that name and version")
)
