package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
)

type Persistence struct {
	pool *pgxpool.Pool
}

func (db *Persistence) Close() { db.pool.Close() }

func (db *Persistence) Pool() *pgxpool.Pool { return db.pool }

func (db *Persistence) HandlerErr() func(error) *domainErr.Error {
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

func NewPersistence(dns string) (*Persistence, error) {
	config, err := pgxpool.ParseConfig(dns)
	if err != nil {
		return nil, err
	}

	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &Persistence{pool: pool}, nil
}
