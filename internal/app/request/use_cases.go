package request

import (
	updatestatus "github.com/MAD-py/pandora-core/internal/app/request/update_status"
	"github.com/MAD-py/pandora-core/internal/validator"
)

// ... Update Status Use Case ...

type UpdateStatusUseCase = updatestatus.UseCase

func NewUpdateStatusUseCase(
	validator validator.Validator, requestRepo RequestUpdateStatusRepository,
) UpdateStatusUseCase {
	return updatestatus.NewUseCase(validator, requestRepo)
}
