package health

import "github.com/MAD-py/pandora-core/internal/app/health/check"

// .. Check Use Case...

type CheckUseCase = check.UseCase

func NewHealthCheckUseCase(database CheckDatabase) CheckUseCase {
	return check.NewUseCase(database)
}
