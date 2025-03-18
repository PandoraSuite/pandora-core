package errors

var (
	ErrAPIKeyExpired          = NewError(CodeValidationError, "api key expired")
	ErrAPIKeyNotFound         = NewError(CodeNotFound, "api key not found")
	ErrAPIKeyGenerationFailed = NewError(CodeInternalError, "api key generation failed")

	ErrClientTypeCannotBeNull = NewError(CodeValidationError, "client type cannot be null")

	ErrServiceNotFound    = NewError(CodeNotFound, "service not found")
	ErrServiceDeprecated  = NewError(CodeValidationError, "service is deprecated")
	ErrServiceDeactivated = NewError(CodeValidationError, "service is deactivated")

	ErrEnvironmentNotFound = NewError(CodeNotFound, "environment not found")

	ErrProjectServiceNotFound = NewError(CodeNotFound, "project service not found")

	ErrProjectStatusCannotBeNull = NewError(CodeValidationError, "project status cannot be null")

	ErrEnvironmentServiceNotFound = NewError(CodeNotFound, "environment service not found")

	ErrNameCannotBeEmpty                     = NewError(CodeValidationError, "name cannot be empty")
	ErrInvalidEmailFormat                    = NewError(CodeValidationError, "invalid email format")
	ErrNoAvailableRequests                   = NewError(CodeValidationError, "no available requests")
	ErrVersionCannotBeEmpty                  = NewError(CodeValidationError, "version cannot be empty")
	ErrMaxRequestExceededForServiceInProyect = NewError(CodeValidationError, "max request exceeded for service in proyect")

	ErrInvalidProjectID  = NewError(CodeValidationError, "project ID must be greater than 0")
	ErrInvalidServiceID  = NewError(CodeValidationError, "service ID must be greater than 0")
	ErrInvalidMaxRequest = NewError(CodeValidationError, "max request must be greater than or equal to 0")
)
