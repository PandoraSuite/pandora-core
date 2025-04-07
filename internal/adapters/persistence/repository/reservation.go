package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ReservationRepository struct {
	pool *pgxpool.Pool

	handlerErr func(error) *errors.Error
}

func (r *ReservationRepository) Save(
	ctx context.Context, Reservation *entities.Reservation,
) *errors.Error {
	//TODO: revisa el orden
	query := `
		INSERT INTO request_log (environment_id, service_id, api_key, request_time, expires_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		Reservation.EnvironmentID,
		Reservation.ServiceID,
		Reservation.APIKey,
		Reservation.RequestTime,
		Reservation.ExpiresAt,
	).Scan(&Reservation.ID)

	return r.handlerErr(err)
}

func NewReservationRepository(
	pool *pgxpool.Pool, handlerErr func(error) *errors.Error,
) *ReservationRepository {
	return &ReservationRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
