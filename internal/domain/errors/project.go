package errors

var (
	ErrProjectNotFound        = NewError(CodeNotFound, "Project not found")
	ErrProjectServiceNotFound = NewError(CodeNotFound, "Project Service not found")

	ErrProjectInvalidStatus                      = NewError(CodeValidationError, "Invalid Project status")
	ErrProjectStatusCannotBeNull                 = NewError(CodeValidationError, "Project status cannot be null")
	ErrProjectServiceAlreadyExists               = NewError(CodeValidationError, "Service is already configured for this Project")
	ErrProjectAlreadyExistsWhitName              = NewError(CodeValidationError, "Project already exists with that name for this Client")
	ErrProjectServiceMaxRequestBelow             = NewError(CodeValidationError, "Mew max request is below the total allocated to Environments")
	ErrProjectServiceNextResetInPast             = NewError(CodeValidationError, "Next reset cannot be set in the past")
	ErrProjectServiceInvalidResetFrequency       = NewError(CodeValidationError, "Invalid Project Service reset frequency")
	ErrProjectServiceResetFrequencyRequired      = NewError(CodeValidationError, "Reset frequency is required when max request is greater than 0")
	ErrProjectServiceNextResetWithInfiniteQuota  = NewError(CodeValidationError, "Next reset cannot be set when max request is unlimited")
	ErrProjectServiceResetFrequencyNotPermitted  = NewError(CodeValidationError, "Reset frequency must be null when max request is 0 (unlimited)")
	ErrProjectServiceFiniteQuotaWithInfiniteEnvs = NewError(CodeValidationError, "Cannot set a finite max request while some Environments have infinite quota")
)
