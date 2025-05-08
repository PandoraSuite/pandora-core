package persistence

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/repository"
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
)

type Repositories interface {
	// ... Helpers ...
	ClosePool()

	// ... Repositories ...
	APIKey() APIKeyRepository
	Client() ClientRepository
	Project() ProjectRepository
	Service() ServiceRepository
	Request() RequestRepository
	Environment() EnvironmentRepository
	Reservation() ReservationRepository
}

var _ Repositories = (*repositories)(nil)

type repositories struct {
	once sync.Once

	pool *pgxpool.Pool

	apiKeyRepo      *repository.APIKeyRepository
	clientRepo      *repository.ClientRepository
	projectRepo     *repository.ProjectRepository
	serviceRepo     *repository.ServiceRepository
	requestRepo     *repository.RequestLogRepository
	environmentRepo *repository.EnvironmentRepository
	reservationRepo *repository.ReservationRepository
}

func (p *repositories) ClosePool() { p.pool.Close() }

func (p *repositories) APIKey() APIKeyRepository {
	p.once.Do(func() {
		p.apiKeyRepo = repository.NewAPIKeyRepository(p.pool, nil)
	})
	return p.apiKeyRepo
}

func (p *repositories) Client() ClientRepository {
	p.once.Do(func() {
		p.clientRepo = repository.NewClientRepository(p.pool, nil)
	})
	return p.clientRepo
}

func (p *repositories) Project() ProjectRepository {
	p.once.Do(func() {
		p.projectRepo = repository.NewProjectRepository(p.pool, nil)
	})
	return p.projectRepo
}

func (p *repositories) Service() ServiceRepository {
	p.once.Do(func() {
		p.serviceRepo = repository.NewServiceRepository(p.pool, nil)
	})
	return p.serviceRepo
}

func (p *repositories) Request() RequestRepository {
	p.once.Do(func() {
		p.requestRepo = repository.NewRequestLogRepository(p.pool, nil)
	})
	return p.requestRepo
}

func (p *repositories) Environment() EnvironmentRepository {
	p.once.Do(func() {
		p.environmentRepo = repository.NewEnvironmentRepository(p.pool, nil)
	})
	return p.environmentRepo
}

func (p *repositories) Reservation() ReservationRepository {
	p.once.Do(func() {
		p.reservationRepo = repository.NewReservationRepository(p.pool, nil)
	})
	return p.reservationRepo
}

func (p *repositories) HandlerErr() func(error) *domainErr.Error {
	uniqueViolation := func(pgErr *pgconn.PgError) *domainErr.Error {
		switch pgErr.ConstraintName {
		case "api_key_key_unique":
			return domainErr.ErrAPIKeyAlreadyExists
		case "client_name_unique":
			return domainErr.ErrClientAlreadyExistsWhitName
		case "client_email_unique":
			return domainErr.ErrClientAlreadyExistsWhitEmail
		case "environment_name_project_id_unique":
			return domainErr.ErrEnvironmentAlreadyExistsWhitName
		case "project_name_client_id_unique":
			return domainErr.ErrProjectAlreadyExistsWhitName
		case "service_name_version_unique":
			return domainErr.ErrServiceAlreadyExistsWhitNameVersion
		default:
			return domainErr.ErrUniqueViolation
		}
	}

	foreignKeyViolation := func(pgErr *pgconn.PgError) *domainErr.Error {
		switch pgErr.ConstraintName {
		case "api_key_environment_id_fk", "request_log_environment_id_fk", "environment_service_environment_id_fk":
			return domainErr.ErrEnvironmentNotFound
		case "environment_project_id_fk", "project_service_project_id_fk":
			return domainErr.ErrProjectNotFound
		case "project_client_id_fk":
			return domainErr.ErrClientNotFound
		case "request_log_service_id_fk", "project_service_service_id_fk", "environment_service_service_id_fk":
			return domainErr.ErrServiceNotFound
		default:
			return domainErr.ErrForeignKeyViolation
		}
	}

	checkViolation := func(pgErr *pgconn.PgError) *domainErr.Error {
		switch pgErr.ConstraintName {
		case "api_key_status_check":
			return domainErr.ErrAPIKeyInvalidStatus
		case "client_type_check":
			return domainErr.ErrClientInvalidType
		case "environment_status_check":
			return domainErr.ErrEnvironmentInvalidStatus
		case "project_service_reset_frequency_check":
			return domainErr.ErrProjectServiceInvalidResetFrequency
		case "request_log_execution_status_check":
			return domainErr.ErrRequestLogInvalidExecutionStatus
		case "service_status_check":
			return domainErr.ErrServiceInvalidStatus
		default:
			return domainErr.ErrRestrictionViolation
		}
	}

	return func(err error) *domainErr.Error {
		if err == nil {
			return nil
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "42P01":
				return domainErr.ErrUndefinedEntity
			case "23505":
				return uniqueViolation(pgErr)
			case "23502":
				return domainErr.ErrNotNullViolation
			case "23503":
				return foreignKeyViolation(pgErr)
			case "23514":
				return checkViolation(pgErr)
			default:
				return domainErr.ErrPersistence
			}
		}
		return domainErr.NewError(
			domainErr.CodeInternalError, "Unknown error", err.Error(),
		)
	}
}

func NewRepositories(dns string) (Repositories, error) {
	config, err := pgxpool.ParseConfig(dns)
	if err != nil {
		return nil, err
	}

	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &repositories{pool: pool}, nil
}
