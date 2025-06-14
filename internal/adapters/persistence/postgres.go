package persistence

import (
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
	if r.driver == nil {
		return
	}

	r.driver.Close()
}

func (r *postgresRepositories) Ping() errors.Error {
	if r.driver == nil {
		return errors.NewInternal(
			"driver is not initialized", nil,
		)
	}
	return r.driver.Ping()
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
