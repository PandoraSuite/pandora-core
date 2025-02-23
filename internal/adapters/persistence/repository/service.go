package repository

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	pool *pgxpool.Pool
}

func (r *ServiceRepository) Save(
	ctx context.Context, service *entities.Service,
) (*entities.Service, error) {
	model := models.ServiceFromEntity(service)
	if err := r.save(ctx, model); err != nil {
		return nil, err
	}

	return model.ToEntity()
}

func (r *ServiceRepository) save(
	ctx context.Context, service *models.Service,
) error {
	if err := service.ValidateModel(); err != nil {
		return err
	}

	query := `
		INSERT INTO service (name, version, status)
		VALUES ($1, $2, $3) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		service.Name,
		service.Version,
		service.Status,
	).Scan(&service.ID, &service.CreatedAt)

	if err != nil {
		return fmt.Errorf("error when inserting the service: %w", err)
	}

	return nil
}

func (r *ServiceRepository) GetByID(
	ctx context.Context, id string,
) (*models.Service, error) {
	query := `
		SELECT id, name, version, status, created_at
		FROM service
		WHERE id = $1;
	`

	var service models.Service
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&service.ID,
		&service.Name,
		&service.Version,
		&service.Status,
		&service.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error obtaining service: %w", err)
	}

	return &service, nil
}
