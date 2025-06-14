package health

import "github.com/MAD-py/pandora-core/internal/app/health/check"

// .. Check Use Case...

type HealthCheckUseCase = check.UseCase

func NewHealthCheckUseCase(database HealthCheckDatabase) HealthCheckUseCase {
	return check.NewUseCase(database)
}
