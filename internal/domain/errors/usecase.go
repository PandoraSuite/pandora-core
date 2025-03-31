package errors

var (
	ErrNameCannotBeEmpty                     = NewError(CodeValidationError, "name cannot be empty")
	ErrInvalidEmailFormat                    = NewError(CodeValidationError, "invalid email format")
	ErrNoAvailableRequests                   = NewError(CodeValidationError, "no available requests")
	ErrVersionCannotBeEmpty                  = NewError(CodeValidationError, "version cannot be empty")
	ErrMaxRequestExceededForServiceInProyect = NewError(CodeValidationError, "max request exceeded for service in proyect")

	ErrInvalidMaxRequest    = NewError(CodeValidationError, "max request must be greater than or equal to 0")
	ErrInvalidEnvironmentID = NewError(CodeValidationError, "environment ID must be greater than 0")

	ErrInfiniteRequestsNotAllowed = NewError(
		CodeValidationError,
		"Infinite requests not allowed",
		"The service cannot have unlimited requests as it has a defined limit in the project",
	)
)
