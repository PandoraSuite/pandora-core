package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func (p *Database) Close() { p.pool.Close() }

func NewDatabase(dns string) (*Database, error) {
	config, err := pgxpool.ParseConfig(dns)
	if err != nil {
		return nil, fmt.Errorf("error al parsear configuración de la BD: %w", err)
	}

	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error al crear el pool de conexiones: %w", err)
	}

	return &Database{pool: pool}, nil
}
