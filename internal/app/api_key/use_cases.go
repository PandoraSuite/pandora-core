package apikey

import (
	"github.com/MAD-py/pandora-core/internal/app/api_key/create"
	"github.com/MAD-py/pandora-core/internal/app/api_key/update"
	validateconsume "github.com/MAD-py/pandora-core/internal/app/api_key/validate_consume"
	validatereservation "github.com/MAD-py/pandora-core/internal/app/api_key/validate_reservation"
	validatereserve "github.com/MAD-py/pandora-core/internal/app/api_key/validate_reserve"
	"github.com/MAD-py/pandora-core/internal/validator"
)

// ... Create Use Case ...

type CreateUseCase = create.UseCase

func NewCreateUseCase(
	validator validator.Validator, repo APIKeyCreateRepository,
) CreateUseCase {
	return create.NewUseCase(validator, repo)
}

// ... Update Use Case ...

type UpdateUseCase = update.UseCase

func NewUpdateUseCase(
	validator validator.Validator, repo APIKeyUpdateRepository,
) UpdateUseCase {
	return update.NewUseCase(validator, repo)
}

// ... Validate And Consume Use Case ...

type ValidateConsumeUseCase = validateconsume.UseCase

func NewValidateConsumeUseCase(
	validator validator.Validator,
	apiKeyRepo APIKeyValidateConsumeRepository,
	requestRepo RequestValidateConsumeRepository,
	serviceRepo ServiceValidateConsumeRepository,
	reservationRepo ReservationValidateConsumeRepository,
	environmentRepo EnvironmentValidateConsumeRepository,
) ValidateConsumeUseCase {
	return validateconsume.NewUseCase(
		validator,
		apiKeyRepo,
		serviceRepo,
		requestRepo,
		reservationRepo,
		environmentRepo,
	)
}

// ... Validate And Reservation Use Case ...

type ValidateReservationUseCase = validatereservation.UseCase

func NewValidateReservationUseCase(
	validator validator.Validator,
	apiKeyRepo APIKeyValidateReservationRepository,
	requestRepo RequestValidateReservationRepository,
	serviceRepo ServiceValidateReservationRepository,
	reservationRepo ReservationValidateReservationRepository,
	environmentRepo EnvironmentValidateReservationRepository,
) ValidateReservationUseCase {
	return validatereservation.NewUseCase(
		validator,
		apiKeyRepo,
		serviceRepo,
		requestRepo,
		reservationRepo,
		environmentRepo,
	)
}

// ... Validate Reservation Use Case ...

type ValidateReserveUseCase = validatereserve.UseCase

func NewValidateReserveUseCase(
	validator validator.Validator,
	requestRepo RequestValidateReserveRepository,
	reservationRepo ReservationValidateReserveRepository,
) ValidateReserveUseCase {
	return validatereserve.NewUseCase(validator, requestRepo, reservationRepo)
}
