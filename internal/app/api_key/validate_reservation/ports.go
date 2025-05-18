package validatereservation

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type APIKeyRepository interface {
	GetByKey(ctx context.Context, key string) (*entities.APIKey, errors.Error)
	UpdateLastUsed(ctx context.Context, key string) errors.Error
}

type EnvironmentRepository interface {
	GetByID(ctx context.Context, id int) (*entities.Environment, errors.Error)
	MissingResourceDiagnosis(ctx context.Context, id int, serviceID int) (bool, bool, errors.Error)
	DecrementAvailableRequest(ctx context.Context, id, serviceID int) (*dto.DecrementAvailableRequest, errors.Error)
}

type ServiceRepository interface {
	GetByNameAndVersion(ctx context.Context, name, version string) (*entities.Service, errors.Error)
}

type RequestRepository interface {
	CreateAsInitialPoint(ctx context.Context, request *entities.RequestLog) errors.Error
}

type ReservationRepository interface {
	Create(ctx context.Context, reservation *entities.Reservation) errors.Error
	CountByEnvironmentAndService(ctx context.Context, environmentID, serviceID int) (int, errors.Error)
}
