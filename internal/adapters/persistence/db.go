package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Persistence struct {
	pool *pgxpool.Pool
}

func (db *Persistence) Close() { db.pool.Close() }

func (db *Persistence) Pool() *pgxpool.Pool { return db.pool }

func NewPersistence(dns string) (*Persistence, error) {
	config, err := pgxpool.ParseConfig(dns)
	if err != nil {
		return nil, fmt.Errorf("error when parsing DB configuration: %w", err)
	}

	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error when creating the connection pool: %w", err)
	}

	return &Persistence{pool: pool}, nil
}
