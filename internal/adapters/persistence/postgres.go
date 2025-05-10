package persistence

import (
	"sync"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/postgres"
)

type postgresRepositories struct {
	driver *postgres.Driver

	once sync.Once

	apiKeyRepo      APIKeyRepository
	clientRepo      ClientRepository
	projectRepo     ProjectRepository
	serviceRepo     ServiceRepository
	requestRepo     RequestRepository
	environmentRepo EnvironmentRepository
	reservationRepo ReservationRepository
}

func (r *postgresRepositories) Close() { r.driver.Close() }

func (r *postgresRepositories) APIKey() APIKeyRepository {
	r.once.Do(func() {
		r.apiKeyRepo = postgres.NewAPIKeyRepository(r.driver)
	})
	return r.apiKeyRepo
}

func (r *postgresRepositories) Client() ClientRepository {
	r.once.Do(func() {
		r.clientRepo = postgres.NewClientRepository(r.driver)
	})
	return r.clientRepo
}

func (r *postgresRepositories) Project() ProjectRepository {
	r.once.Do(func() {
		r.projectRepo = postgres.NewProjectRepository(r.driver)
	})
	return r.projectRepo
}

func (r *postgresRepositories) Service() ServiceRepository {
	r.once.Do(func() {
		r.serviceRepo = postgres.NewServiceRepository(r.driver)
	})
	return r.serviceRepo
}

func (r *postgresRepositories) Request() RequestRepository {
	r.once.Do(func() {
		r.requestRepo = postgres.NewRequestLogRepository(r.driver)
	})
	return r.requestRepo
}

func (r *postgresRepositories) Environment() EnvironmentRepository {
	r.once.Do(func() {
		r.environmentRepo = postgres.NewEnvironmentRepository(r.driver)
	})
	return r.environmentRepo
}

func (r *postgresRepositories) Reservation() ReservationRepository {
	r.once.Do(func() {
		r.reservationRepo = postgres.NewReservationRepository(r.driver)
	})
	return r.reservationRepo
}
