package persistence

import (
	"context"
	"time"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/postgres"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports"
)

type postgresRepositories struct {
	driver *postgres.Driver

	apiKeyRepo      ports.APIKeyRepository
	clientRepo      ports.ClientRepository
	projectRepo     ports.ProjectRepository
	serviceRepo     ports.ServiceRepository
	requestRepo     ports.RequestRepository
	environmentRepo ports.EnvironmentRepository
	reservationRepo ports.ReservationRepository
}

func (r *postgresRepositories) Close() {
	if pool := r.driver.Pool(); pool != nil {
		pool.Close()
	}
}

func (r *postgresRepositories) Ping() errors.Error {
	pool := r.driver.Pool()

	if pool == nil {
		return errors.NewInternal("pool is not initialized", nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return errors.NewInternal(
			"failed to ping Postgres database",
			err,
		)
	}
	return nil
}

func (r *postgresRepositories) Latency() (int64, errors.Error) {
	pool := r.driver.Pool()
	if pool == nil {
		return 0, errors.NewInternal("pool is not initialized", nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	err := pool.QueryRow(ctx, "SELECT 1").Scan(new(int))
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return 0, errors.NewInternal("failed to measure DB latency", err)
	}

	return latency, nil
}

func (r *postgresRepositories) APIKey() ports.APIKeyRepository {
	if r.apiKeyRepo == nil {
		r.apiKeyRepo = postgres.NewAPIKeyRepository(r.driver)
	}
	return r.apiKeyRepo
}

func (r *postgresRepositories) Client() ports.ClientRepository {
	if r.clientRepo == nil {
		r.clientRepo = postgres.NewClientRepository(r.driver)
	}
	return r.clientRepo
}

func (r *postgresRepositories) Project() ports.ProjectRepository {
	if r.projectRepo == nil {
		r.projectRepo = postgres.NewProjectRepository(r.driver)
	}
	return r.projectRepo
}

func (r *postgresRepositories) Service() ports.ServiceRepository {
	if r.serviceRepo == nil {
		r.serviceRepo = postgres.NewServiceRepository(r.driver)
	}
	return r.serviceRepo
}

func (r *postgresRepositories) Request() ports.RequestRepository {
	if r.requestRepo == nil {
		r.requestRepo = postgres.NewRequestRepository(r.driver)
	}
	return r.requestRepo
}

func (r *postgresRepositories) Environment() ports.EnvironmentRepository {
	if r.environmentRepo == nil {
		r.environmentRepo = postgres.NewEnvironmentRepository(r.driver)
	}
	return r.environmentRepo
}

func (r *postgresRepositories) Reservation() ports.ReservationRepository {
	if r.reservationRepo == nil {
		r.reservationRepo = postgres.NewReservationRepository(r.driver)
	}
	return r.reservationRepo
}
