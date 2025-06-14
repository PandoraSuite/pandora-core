package check

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type UseCase interface {
	Execute() *dto.HealthCheckResponse
}

type useCase struct {
	database Database
}

func (uc *useCase) Execute() *dto.HealthCheckResponse {
	databaseCheck := uc.checkDatabase()
	return &dto.HealthCheckResponse{
		Status:    databaseCheck.Status,
		Timestamp: time.Now(),
		Check: &dto.CheckResponse{
			Database: databaseCheck,
		},
	}
}

func (uc *useCase) checkDatabase() *dto.CheckStatusResponse {
	if err := uc.database.Ping(); err != nil {
		return &dto.CheckStatusResponse{
			Status:  enums.HealthStatusDown,
			Message: err.Error(),
		}
	}

	latency, err := uc.database.Latency()
	if err != nil {
		return &dto.CheckStatusResponse{
			Status:  enums.HealthStatusDegraded,
			Message: err.Error(),
		}
	}

	return &dto.CheckStatusResponse{
		Status:  enums.HealthStatusOK,
		Message: "database is reachable",
		Latency: latency,
	}
}

func NewUseCase(database Database) UseCase {
	return &useCase{
		database: database,
	}
}
