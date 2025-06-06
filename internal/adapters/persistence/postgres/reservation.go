package postgres

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ReservationRepository struct {
	*Driver

	tableName string
}

func (r *ReservationRepository) Create(
	ctx context.Context, Reservation *entities.Reservation,
) errors.Error {
	query := `
		INSERT INTO reservation (environment_id, service_id, api_key, start_request_id, request_time, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		Reservation.EnvironmentID,
		Reservation.ServiceID,
		Reservation.APIKey,
		Reservation.StartRequestID,
		Reservation.RequestTime,
		Reservation.ExpiresAt,
	).Scan(&Reservation.ID)

	return r.errorMapper(err, r.tableName)
}

func (r *ReservationRepository) GetByID(
	ctx context.Context, id string,
) (*entities.Reservation, errors.Error) {
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
		return nil, r.errorMapper(err, r.tableName)
	}

	return reservation, nil
}

func (r *ReservationRepository) GetByIDWithDetails(
	ctx context.Context, id string,
) (*dto.ReservationWithDetails, errors.Error) {
	query := `
		SELECT r.id, r.start_request_id, r.api_key, 
		s.id, s.name, s.version, s.status, 
		e.id, e.name, e.status
		FROM reservation r
		INNER JOIN environment_service es ON r.environment_id = es.environment_id 
		AND r.service_id = es.service_id
		INNER JOIN environment e ON es.environment_id = e.id
		INNER JOIN service s ON es.service_id = s.id
		WHERE r.id = $1;
	`

	reservationFlow := new(dto.ReservationWithDetails)
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&reservationFlow.ID,
		&reservationFlow.StartRequestID,
		&reservationFlow.APIKey,
		&reservationFlow.ServiceID,
		&reservationFlow.ServiceName,
		&reservationFlow.ServiceVersion,
		&reservationFlow.ServiceStatus,
		&reservationFlow.EnvironmentID,
		&reservationFlow.EnvironmentName,
		&reservationFlow.EnvironmentStatus,
	)
	if err != nil {
		return nil, r.errorMapper(err, r.tableName)
	}

	return reservationFlow, nil
}

func (r *ReservationRepository) CountByEnvironmentAndService(
	ctx context.Context, environment_id, service_id int,
) (int, errors.Error) {
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
		return 0, r.errorMapper(err, r.tableName)
	}

	return currentReservations, nil
}

func (r *ReservationRepository) Delete(
	ctx context.Context, id string,
) errors.Error {
	query := `
		DELETE FROM reservation
		WHERE id = $1;
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return r.errorMapper(err, r.tableName)
	}

	if result.RowsAffected() == 0 {
		return r.entityNotFoundError(r.tableName, map[string]any{"id": id})
	}

	return nil
}

func NewReservationRepository(driver *Driver) *ReservationRepository {
	return &ReservationRepository{Driver: driver, tableName: "reservation"}
}
