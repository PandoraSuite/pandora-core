package errors

var (
	ErrReservationNotFound         = NewError(CodeNotFound, "Reservation not found")
	ErrReservationGenerationFailed = NewError(CodeInternalError, "Reservation generation failed")
)
