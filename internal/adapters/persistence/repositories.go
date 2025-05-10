package persistence

import (
	"context"
	"time"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/errors"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type DriverType string

const (
	PostgresDriver DriverType = "postgres"
)

type Repositories interface {
	// ... Helpers ...
	Close()

	// ... Repositories ...
	APIKey() APIKeyRepository
	Client() ClientRepository
	Project() ProjectRepository
	Service() ServiceRepository
	Request() RequestRepository
	Environment() EnvironmentRepository
	Reservation() ReservationRepository
}

type APIKeyRepository interface {
	// ... Exists ...
	Exists(ctx context.Context, key string) (bool, errors.Error)

	// ... Get ...
	GetByID(ctx context.Context, id int) (*entities.APIKey, errors.Error)
	GetByKey(ctx context.Context, key string) (*entities.APIKey, errors.Error)

	// ... List ...
	ListByEnvironment(ctx context.Context, environmentID int) ([]*entities.APIKey, errors.Error)

	// ... Create ...
	Create(ctx context.Context, apiKey *entities.APIKey) errors.Error

	// ... Update ...
	Update(ctx context.Context, id int, update *dto.APIKeyUpdate) (*entities.APIKey, errors.Error)
	UpdateStatus(ctx context.Context, id int, status enums.APIKeyStatus) errors.Error
	UpdateLastUsed(ctx context.Context, key string) errors.Error
}

type ClientRepository interface {
	// ... Exists ...
	Exists(ctx context.Context, id int) (bool, errors.Error)

	// ... Get ...
	GetByID(ctx context.Context, id int) (*entities.Client, errors.Error)

	// ... List ...
	List(ctx context.Context, filter *dto.ClientFilter) ([]*entities.Client, errors.Error)

	// ... Create ...
	Create(ctx context.Context, client *entities.Client) errors.Error

	// ... Update ...
	Update(ctx context.Context, id int, update *dto.ClientUpdate) (*entities.Client, errors.Error)
}

type EnvironmentRepository interface {
	// ... Exists ...
	Exists(ctx context.Context, id int) (bool, errors.Error)
	IsActive(ctx context.Context, id int) (bool, errors.Error)
	ExistsServiceIn(ctx context.Context, id, serviceID int) (bool, errors.Error)
	MissingResourceDiagnosis(ctx context.Context, id int, serviceID int) (bool, bool, errors.Error)
	ExistsServiceWithInfiniteMaxRequest(ctx context.Context, projectID, serviceID int) (bool, errors.Error)

	// ... Get ...
	GetByID(ctx context.Context, id int) (*entities.Environment, errors.Error)
	GetServiceByID(ctx context.Context, id, serviceID int) (*entities.EnvironmentService, errors.Error)
	GetProjectServiceQuotaUsage(ctx context.Context, id, serviceID int) (*dto.QuotaUsage, errors.Error)

	// ... List ...
	ListByProject(ctx context.Context, projectID int) ([]*entities.Environment, errors.Error)

	// ... Create ...
	Create(ctx context.Context, environment *entities.Environment) errors.Error
	AddService(ctx context.Context, id int, service *entities.EnvironmentService) errors.Error

	// ... Update ...
	Update(ctx context.Context, id int, update *dto.EnvironmentUpdate) (*entities.Environment, errors.Error)
	UpdateStatus(ctx context.Context, id int, status enums.EnvironmentStatus) errors.Error
	UpdateService(ctx context.Context, id, serviceID int, update *dto.EnvironmentServiceUpdate) (*entities.EnvironmentService, errors.Error)
	ResetAvailableRequests(ctx context.Context, id, serviceID int) (*entities.EnvironmentService, errors.Error)
	IncreaseAvailableRequest(ctx context.Context, id, serviceID int) errors.Error
	DecrementAvailableRequest(ctx context.Context, id, serviceID int) (*dto.DecrementAvailableRequest, errors.Error)

	// ... Delete ...
	RemoveService(ctx context.Context, id, serviceID int) (int64, errors.Error)
	RemoveServiceFromProjectEnvironments(ctx context.Context, projectID, serviceID int) (int64, errors.Error)
}

type ProjectRepository interface {
	// ... Exists ...
	Exists(ctx context.Context, id int) (bool, errors.Error)
	ExistsServiceIn(ctx context.Context, serviceID int) (bool, errors.Error)

	// ... Get ...
	GetByID(ctx context.Context, id int) (*entities.Project, errors.Error)
	GetServiceByID(ctx context.Context, id, serviceID int) (*entities.ProjectService, errors.Error)
	GetProjectServiceQuotaUsage(ctx context.Context, id, serviceID int) (*dto.QuotaUsage, errors.Error)

	// ... List ...
	List(ctx context.Context) ([]*entities.Project, errors.Error)
	ListByClient(ctx context.Context, clientID int) ([]*entities.Project, errors.Error)

	// ... Create ...
	Create(ctx context.Context, project *entities.Project) errors.Error
	AddService(ctx context.Context, id int, service *entities.ProjectService) errors.Error

	// ... Update ...
	Update(ctx context.Context, id int, update *dto.ProjectUpdate) (*entities.Project, errors.Error)
	UpdateStatus(ctx context.Context, id int, status enums.ProjectStatus) errors.Error
	UpdateService(ctx context.Context, id, serviceID int, update *dto.ProjectServiceUpdate) (*entities.ProjectService, errors.Error)
	ResetProjectServiceUsage(ctx context.Context, id, serviceID int, nextReset time.Time) ([]*dto.EnvironmentServiceReset, errors.Error)
	ResetAvailableRequestsForEnvsService(ctx context.Context, id, serviceID int) ([]*dto.EnvironmentServiceReset, errors.Error)

	// ... Delete ...
	RemoveService(ctx context.Context, id, serviceID int) (int64, errors.Error)
}

type RequestRepository interface {
	// ... Create ...
	Create(ctx context.Context, requestLog *entities.RequestLog) errors.Error
	CreateAsInitialPoint(ctx context.Context, requestLog *entities.RequestLog) errors.Error

	// ... Update ...
	UpdateExecutionStatus(ctx context.Context, id string, executionStatus enums.RequestLogExecutionStatus) errors.Error

	// ... Delete ...
	DeleteByService(ctx context.Context, serviceID int) errors.Error
}

type ReservationRepository interface {
	// ... Get ...
	GetByID(ctx context.Context, id string) (*entities.Reservation, errors.Error)
	GetByIDWithDetails(ctx context.Context, id string) (*dto.ReservationWithDetails, errors.Error)
	CountByEnvironmentAndService(ctx context.Context, environment_id, service_id int) (int, errors.Error)

	// ... Create ...
	Create(ctx context.Context, Reservation *entities.Reservation) errors.Error

	// ... Delete ...
	Delete(ctx context.Context, id string) errors.Error
}

type ServiceRepository interface {
	// ... Get ...
	GetByNameAndVersion(ctx context.Context, name, version string) (*entities.Service, errors.Error)

	// ... List ...
	List(ctx context.Context, filter *dto.ServiceFilter) ([]*entities.Service, errors.Error)

	// ... Create ...
	Create(ctx context.Context, service *entities.Service) errors.Error

	// ... Update ...
	UpdateStatus(ctx context.Context, id int, status enums.ServiceStatus) (*entities.Service, errors.Error)

	// ... Delete ...
	Delete(ctx context.Context, id int) errors.Error
}
