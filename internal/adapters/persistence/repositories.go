package persistence

import "github.com/MAD-py/pandora-core/internal/ports"

type DriverType string

const (
	PostgresDriver DriverType = "postgres"
)

type Repositories interface {
	// ... Helpers ...
	Close()

	// ... Repositories ...
	APIKey() ports.APIKeyRepository
	Client() ports.ClientRepository
	Project() ports.ProjectRepository
	Service() ports.ServiceRepository
	Request() ports.RequestRepository
	Environment() ports.EnvironmentRepository
	Reservation() ports.ReservationRepository
}
