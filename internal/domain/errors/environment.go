package errors

var (
	ErrEnvironmentNotFound        = NewError(CodeNotFound, "Environment not found")
	ErrEnvironmentServiceNotFound = NewError(CodeNotFound, "Environment Service not found")

	ErrEnvironmentInvalidStatus                     = NewError(CodeValidationError, "Invalid Environment status")
	ErrEnvironmentServiceAlreadyExists              = NewError(CodeValidationError, "Service is already configured for this Environment")
	ErrEnvironmentAlreadyExistsWhitName             = NewError(CodeValidationError, "Environment already exists with that name for this Project")
	ErrEnvironmentServiceAvailableRequestExceedsMax = NewError(CodeValidationError, "Available request cannot be greater than max request")
	ErrEnvironmentServiceAvailableRequestNotAllowed = NewError(CodeValidationError, "Available request cannot be set when max request is 0 (unlimited)")
)
