package repository

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnvironmentRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) error
}

func (r *EnvironmentRepository) FindByID(
	ctx context.Context, id int,
) (*entities.Environment, error) {
	query := "SELECT * FROM environment WHERE id = $1;"

	var environment models.Environment
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&environment.ID,
		&environment.ProjectID,
		&environment.Name,
		&environment.Status,
		&environment.CreatedAt,
	)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	return environment.ToEntity()
}

func (r *EnvironmentRepository) FindByProject(
	ctx context.Context, projectID int,
) ([]*entities.Environment, error) {
	query := "SELECT * FROM environment WHERE project_id = $1;"

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	defer rows.Close()

	var environments []*models.Environment
	for rows.Next() {
		environment := new(models.Environment)

		err = rows.Scan(
			&environment.ID,
			&environment.ProjectID,
			&environment.Name,
			&environment.Status,
			&environment.CreatedAt,
		)
		if err != nil {
			return nil, r.handlerErr(err)
		}

		environments = append(environments, environment)
	}

	if err := rows.Err(); err != nil {
		return nil, r.handlerErr(err)
	}

	return models.EnvironmentToEntity(environments)
}

func (r *EnvironmentRepository) Save(
	ctx context.Context, environment *entities.Environment,
) error {
	model := models.EnvironmentFromEntity(environment)
	if err := r.save(ctx, &model); err != nil {
		return err
	}

	environment.ID = model.EntityID()
	environment.CreatedAt = model.EntityCreatedAt()
	return nil
}

func (r *EnvironmentRepository) save(
	ctx context.Context, environment *models.Environment,
) error {
	if err := environment.ValidateModel(); err != nil {
		return err
	}

	query := `
		INSERT INTO environment (project_id, name, status)
		VALUES ($1, $2, $3) RETURNING id, created_at;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		environment.ProjectID,
		environment.Name,
		environment.Status,
	).Scan(&environment.ID, &environment.CreatedAt)

	if err != nil {
		return r.handlerErr(err)
	}

	return nil
}

func NewEnvironmentRepository(
	pool *pgxpool.Pool, handlerErr func(error) error,
) *EnvironmentRepository {
	return &EnvironmentRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
