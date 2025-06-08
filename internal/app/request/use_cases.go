package request

import (
	updateexecutionstatus "github.com/MAD-py/pandora-core/internal/app/request/update_execution_status"
	"github.com/MAD-py/pandora-core/internal/validator"
)

// ... Update Execution Status Use Case ...

type UpdateExecutionStatusUseCase = updateexecutionstatus.UseCase

func NewUpdateExecutionStatusUseCase(
	validator validator.Validator, requestRepo RequestUpdateExecutionStatusRepository,
) UpdateExecutionStatusUseCase {
	return updateexecutionstatus.NewUseCase(validator, requestRepo)
}
