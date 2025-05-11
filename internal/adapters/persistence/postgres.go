package persistence

import (
	"sync"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/postgres"
	"github.com/MAD-py/pandora-core/internal/ports"
)

type postgresRepositories struct {
	driver *postgres.Driver

	once sync.Once

	apiKeyRepo      ports.APIKeyRepository
	clientRepo      ports.ClientRepository
	projectRepo     ports.ProjectRepository
	serviceRepo     ports.ServiceRepository
	requestRepo     ports.RequestRepository
	environmentRepo ports.EnvironmentRepository
	reservationRepo ports.ReservationRepository
}

func (r *postgresRepositories) Close() { r.driver.Close() }

func (r *postgresRepositories) APIKey() ports.APIKeyRepository {
	r.once.Do(func() {
		r.apiKeyRepo = postgres.NewAPIKeyRepository(r.driver)
	})
	return r.apiKeyRepo
}

func (r *postgresRepositories) Client() ports.ClientRepository {
	r.once.Do(func() {
		r.clientRepo = postgres.NewClientRepository(r.driver)
	})
	return r.clientRepo
}

func (r *postgresRepositories) Project() ports.ProjectRepository {
	r.once.Do(func() {
		r.projectRepo = postgres.NewProjectRepository(r.driver)
	})
	return r.projectRepo
}

func (r *postgresRepositories) Service() ports.ServiceRepository {
	r.once.Do(func() {
		r.serviceRepo = postgres.NewServiceRepository(r.driver)
	})
	return r.serviceRepo
}

func (r *postgresRepositories) Request() ports.RequestRepository {
	r.once.Do(func() {
		r.requestRepo = postgres.NewRequestLogRepository(r.driver)
	})
	return r.requestRepo
}

func (r *postgresRepositories) Environment() ports.EnvironmentRepository {
	r.once.Do(func() {
		r.environmentRepo = postgres.NewEnvironmentRepository(r.driver)
	})
	return r.environmentRepo
}

func (r *postgresRepositories) Reservation() ports.ReservationRepository {
	r.once.Do(func() {
		r.reservationRepo = postgres.NewReservationRepository(r.driver)
	})
	return r.reservationRepo
}
