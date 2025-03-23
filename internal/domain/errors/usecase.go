package errors

var (
	ErrAPIKeyExpired          = NewError(CodeValidationError, "api key expired")
	ErrAPIKeyNotFound         = NewError(CodeNotFound, "api key not found")
	ErrAPIKeyInvalidStatus    = NewError(CodeValidationError, "invalid API key status")
	ErrAPIKeyGenerationFailed = NewError(CodeInternalError, "api key generation failed")

	ErrClientNotFound         = NewError(CodeNotFound, "client not found")
	ErrClientInvalidType      = NewError(CodeValidationError, "invalid client type")
	ErrClientTypeCannotBeNull = NewError(CodeValidationError, "client type cannot be null")

	ErrEnvironmentNotFound      = NewError(CodeNotFound, "environment not found")
	ErrEnvironmentInvalidStatus = NewError(CodeValidationError, "invalid environment status")

	ErrProjectServiceNotFound      = NewError(CodeNotFound, "project service not found")
	ErrProjectServiceInvalidStatus = NewError(CodeValidationError, "invalid project service status")

	ErrProjectNotFound           = NewError(CodeNotFound, "project not found")
	ErrProjectInvalidStatus      = NewError(CodeValidationError, "invalid project status")
	ErrProjectStatusCannotBeNull = NewError(CodeValidationError, "project status cannot be null")

	ErrServiceNotFound      = NewError(CodeNotFound, "service not found")
	ErrServiceDeprecated    = NewError(CodeValidationError, "service is deprecated")
	ErrServiceDeactivated   = NewError(CodeValidationError, "service is deactivated")
	ErrServiceInvalidStatus = NewError(CodeValidationError, "invalid service status")

	ErrEnvironmentServiceNotFound = NewError(CodeNotFound, "environment service not found")

	ErrNameCannotBeEmpty                     = NewError(CodeValidationError, "name cannot be empty")
	ErrInvalidEmailFormat                    = NewError(CodeValidationError, "invalid email format")
	ErrNoAvailableRequests                   = NewError(CodeValidationError, "no available requests")
	ErrVersionCannotBeEmpty                  = NewError(CodeValidationError, "version cannot be empty")
	ErrMaxRequestExceededForServiceInProyect = NewError(CodeValidationError, "max request exceeded for service in proyect")

	ErrInvalidClientID   = NewError(CodeValidationError, "client ID must be greater than 0")
	ErrInvalidProjectID  = NewError(CodeValidationError, "project ID must be greater than 0")
	ErrInvalidServiceID  = NewError(CodeValidationError, "service ID must be greater than 0")
	ErrInvalidMaxRequest = NewError(CodeValidationError, "max request must be greater than or equal to 0")
)
