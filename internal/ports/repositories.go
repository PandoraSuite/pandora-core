package ports

import (
	"context"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type APIKeyRepository interface {
	// ... Exists ...
	Exists(ctx context.Context, key string) (bool, errors.PersistenceError)

	// ... Get ...
	GetByID(ctx context.Context, id int) (*entities.APIKey, errors.PersistenceError)
	GetByKey(ctx context.Context, key string) (*entities.APIKey, errors.PersistenceError)

	// ... List ...
	ListByEnvironment(ctx context.Context, environmentID int) ([]*entities.APIKey, errors.PersistenceError)

	// ... Create ...
	Create(ctx context.Context, apiKey *entities.APIKey) errors.PersistenceError

	// ... Update ...
	Update(ctx context.Context, id int, update *dto.APIKeyUpdate) (*entities.APIKey, errors.PersistenceError)
	UpdateStatus(ctx context.Context, id int, status enums.APIKeyStatus) errors.PersistenceError
	UpdateLastUsed(ctx context.Context, key string) errors.PersistenceError
}

type ClientRepository interface {
	// ... Exists ...
	Exists(ctx context.Context, id int) (bool, errors.PersistenceError)

	// ... Get ...
	GetByID(ctx context.Context, id int) (*entities.Client, errors.PersistenceError)

	// ... List ...
	List(ctx context.Context, filter *dto.ClientFilter) ([]*entities.Client, errors.PersistenceError)

	// ... Create ...
	Create(ctx context.Context, client *entities.Client) errors.PersistenceError

	// ... Update ...
	Update(ctx context.Context, id int, update *dto.ClientUpdate) (*entities.Client, errors.PersistenceError)
}

type EnvironmentRepository interface {
	// ... Exists ...
	Exists(ctx context.Context, id int) (bool, errors.PersistenceError)
	IsActive(ctx context.Context, id int) (bool, errors.PersistenceError)
	ExistsServiceIn(ctx context.Context, id, serviceID int) (bool, errors.PersistenceError)
	MissingResourceDiagnosis(ctx context.Context, id int, serviceID int) (bool, bool, errors.PersistenceError)
	ExistsServiceWithInfiniteMaxRequest(ctx context.Context, projectID, serviceID int) (bool, errors.PersistenceError)

	// ... Get ...
	GetByID(ctx context.Context, id int) (*entities.Environment, errors.PersistenceError)
	GetServiceByID(ctx context.Context, id, serviceID int) (*entities.EnvironmentService, errors.PersistenceError)
	GetProjectServiceQuotaUsage(ctx context.Context, id, serviceID int) (*dto.QuotaUsage, errors.PersistenceError)

	// ... List ...
	ListByProject(ctx context.Context, projectID int) ([]*entities.Environment, errors.PersistenceError)

	// ... Create ...
	Create(ctx context.Context, environment *entities.Environment) errors.PersistenceError
	AddService(ctx context.Context, id int, service *entities.EnvironmentService) errors.PersistenceError

	// ... Update ...
	Update(ctx context.Context, id int, update *dto.EnvironmentUpdate) (*entities.Environment, errors.PersistenceError)
	UpdateStatus(ctx context.Context, id int, status enums.EnvironmentStatus) errors.PersistenceError
	UpdateService(ctx context.Context, id, serviceID int, update *dto.EnvironmentServiceUpdate) (*entities.EnvironmentService, errors.PersistenceError)
	ResetAvailableRequests(ctx context.Context, id, serviceID int) (*entities.EnvironmentService, errors.PersistenceError)
	IncreaseAvailableRequest(ctx context.Context, id, serviceID int) errors.PersistenceError
	DecrementAvailableRequest(ctx context.Context, id, serviceID int) (*dto.DecrementAvailableRequest, errors.PersistenceError)

	// ... Delete ...
	RemoveService(ctx context.Context, id, serviceID int) (int64, errors.PersistenceError)
	RemoveServiceFromProjectEnvironments(ctx context.Context, projectID, serviceID int) (int64, errors.PersistenceError)
}

type ProjectRepository interface {
	// ... Exists ...
	Exists(ctx context.Context, id int) (bool, errors.PersistenceError)
	ExistsServiceIn(ctx context.Context, serviceID int) (bool, errors.PersistenceError)

	// ... Get ...
	GetByID(ctx context.Context, id int) (*entities.Project, errors.PersistenceError)
	GetServiceByID(ctx context.Context, id, serviceID int) (*entities.ProjectService, errors.PersistenceError)
	GetProjectServiceQuotaUsage(ctx context.Context, id, serviceID int) (*dto.QuotaUsage, errors.PersistenceError)

	// ... List ...
	List(ctx context.Context) ([]*entities.Project, errors.PersistenceError)
	ListByClient(ctx context.Context, clientID int) ([]*entities.Project, errors.PersistenceError)

	// ... Create ...
	Create(ctx context.Context, project *entities.Project) errors.PersistenceError
	AddService(ctx context.Context, id int, service *entities.ProjectService) errors.PersistenceError

	// ... Update ...
	Update(ctx context.Context, id int, update *dto.ProjectUpdate) (*entities.Project, errors.PersistenceError)
	UpdateStatus(ctx context.Context, id int, status enums.ProjectStatus) errors.PersistenceError
	UpdateService(ctx context.Context, id, serviceID int, update *dto.ProjectServiceUpdate) (*entities.ProjectService, errors.PersistenceError)
	ResetProjectServiceUsage(ctx context.Context, id, serviceID int, nextReset time.Time) ([]*dto.EnvironmentServiceReset, errors.PersistenceError)
	ResetAvailableRequestsForEnvsService(ctx context.Context, id, serviceID int) ([]*dto.EnvironmentServiceReset, errors.PersistenceError)

	// ... Delete ...
	RemoveService(ctx context.Context, id, serviceID int) (int64, errors.PersistenceError)
}

type RequestRepository interface {
	// ... Create ...
	Create(ctx context.Context, requestLog *entities.RequestLog) errors.PersistenceError
	CreateAsInitialPoint(ctx context.Context, requestLog *entities.RequestLog) errors.PersistenceError

	// ... Update ...
	UpdateExecutionStatus(ctx context.Context, id string, executionStatus enums.RequestLogExecutionStatus) errors.PersistenceError

	// ... Delete ...
	DeleteByService(ctx context.Context, serviceID int) errors.PersistenceError
}

type ReservationRepository interface {
	// ... Get ...
	GetByID(ctx context.Context, id string) (*entities.Reservation, errors.PersistenceError)
	GetByIDWithDetails(ctx context.Context, id string) (*dto.ReservationWithDetails, errors.PersistenceError)
	CountByEnvironmentAndService(ctx context.Context, environment_id, service_id int) (int, errors.PersistenceError)

	// ... Create ...
	Create(ctx context.Context, Reservation *entities.Reservation) errors.PersistenceError

	// ... Delete ...
	Delete(ctx context.Context, id string) errors.PersistenceError
}

type ServiceRepository interface {
	// ... Get ...
	GetByNameAndVersion(ctx context.Context, name, version string) (*entities.Service, errors.PersistenceError)

	// ... List ...
	List(ctx context.Context, filter *dto.ServiceFilter) ([]*entities.Service, errors.PersistenceError)

	// ... Create ...
	Create(ctx context.Context, service *entities.Service) errors.PersistenceError

	// ... Update ...
	UpdateStatus(ctx context.Context, id int, status enums.ServiceStatus) (*entities.Service, errors.PersistenceError)

	// ... Delete ...
	Delete(ctx context.Context, id int) errors.PersistenceError
}

type CredentialsRepository interface {
	// ... Get ...
	GetByUsername(ctx context.Context, username string) (*entities.Credentials, errors.PersistenceError)

	// ... Update ...
	ChangePassword(ctx context.Context, credentials *entities.Credentials) errors.PersistenceError
}
