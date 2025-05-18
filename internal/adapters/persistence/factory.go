package persistence

import (
	"github.com/MAD-py/pandora-core/internal/adapters/persistence/postgres"
)

func NewRepositories(driver DriverType, dns string) Repositories {
	switch driver {
	case PostgresDriver:
		return &postgresRepositories{driver: postgres.NewDriver(dns)}
	default:
		panic("unsupported driver type " + string(driver))
	}
}
