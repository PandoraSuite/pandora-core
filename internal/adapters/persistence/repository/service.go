package repository

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	pool *pgxpool.Pool
}

func (r *ServiceRepository) FindByNameAndVersion(
	ctx context.Context, name, version string,
) (*entities.Service, error) {
	query := `SELECT * FROM service WHERE name = $1 AND version = $2;`

	var service models.Service
	err := r.pool.QueryRow(ctx, query, name, version).Scan(
		&service.ID,
		&service.Name,
		&service.Version,
		&service.Status,
		&service.CreatedAt,
	)
	if err != nil {
		return nil, persistence.ConvertPgxError(err)
	}

	return service.ToEntity()
}

func (r *ServiceRepository) FindActiveServices(ctx context.Context) ([]*entities.Service, error) {
	query := "SELECT * FROM service WHERE status = 'active';"
	return r.find(ctx, query)
}

func (r *ServiceRepository) FindAll(ctx context.Context) ([]*entities.Service, error) {
	query := "SELECT * FROM service;"
	return r.find(ctx, query)
}

func (r *ServiceRepository) find(
	ctx context.Context, query string,
) ([]*entities.Service, error) {
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, persistence.ConvertPgxError(err)
	}

	defer rows.Close()

	var services []*models.Service
	for rows.Next() {
		service := new(models.Service)

		err = rows.Scan(
			&service.ID,
			&service.Name,
			&service.Version,
			&service.Status,
			&service.CreatedAt,
		)
		if err != nil {
			return nil, persistence.ConvertPgxError(err)
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, persistence.ConvertPgxError(err)
	}

	return models.ServicesToEntity(services)
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
		return persistence.ConvertPgxError(err)
	}

	return nil
}

func NewServiceRepository(pool *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{pool: pool}
}
