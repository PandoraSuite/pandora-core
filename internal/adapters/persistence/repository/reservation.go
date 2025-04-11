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
	query := `
		INSERT INTO reservation (environment_id, service_id, api_key, request_time, expires_at)
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

func (r *ReservationRepository) FindByID(
	ctx context.Context, id string,
) (*entities.Reservation, *errors.Error) {
	query := `
		SELECT *
		FROM reservation
		WHERE id = $1;
	`

	reservation := new(entities.Reservation)
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&reservation.ID,
		&reservation.EnvironmentID,
		&reservation.ServiceID,
		&reservation.APIKey,
		&reservation.RequestTime,
		&reservation.ExpiresAt,
	)
	if err != nil {
		return nil, r.handlerErr(err)
	}

	return reservation, nil
}

func (r *ReservationRepository) CountByEnvironmentAndService(
	ctx context.Context, environment_id, service_id int,
) (int, *errors.Error) {
	query := `
		SELECT count(*)
		FROM reservation
		WHERE environment_id = $1
		AND service_id = $2
	`

	var currentReservations int
	err := r.pool.QueryRow(
		ctx,
		query,
		environment_id,
		service_id).Scan(
		&currentReservations)
	if err != nil {
		return 0, r.handlerErr(err)
	}

	return currentReservations, nil
}

func (r *ReservationRepository) RemoveReservation(
	ctx context.Context, id string,
) (int64, *errors.Error) {
	query := `
		DELETE FROM reservation
		WHERE id = $1;
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return 0, r.handlerErr(err)
	}

	if result.RowsAffected() == 0 {
		return 0, errors.ErrReservationNotFound
	}

	return result.RowsAffected(), nil
}
func NewReservationRepository(
	pool *pgxpool.Pool, handlerErr func(error) *errors.Error,
) *ReservationRepository {
	return &ReservationRepository{
		pool:       pool,
		handlerErr: handlerErr,
	}
}
