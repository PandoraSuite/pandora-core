package apikey

import (
	"github.com/MAD-py/pandora-core/internal/app/api_key/create"
	"github.com/MAD-py/pandora-core/internal/app/api_key/update"
	validateconsume "github.com/MAD-py/pandora-core/internal/app/api_key/validate_consume"
	validatereservation "github.com/MAD-py/pandora-core/internal/app/api_key/validate_reservation"
	validatereserve "github.com/MAD-py/pandora-core/internal/app/api_key/validate_reserve"
)

// ... Create Use Case ...

type APIKeyCreateRepository = create.APIKeyRepository

// ... Update Use Case ...

type APIKeyUpdateRepository = update.APIKeyRepository

// ... Validate And Consume Use Case ...

type APIKeyValidateConsumeRepository = validateconsume.APIKeyRepository
type RequestValidateConsumeRepository = validateconsume.RequestRepository
type ServiceValidateConsumeRepository = validateconsume.ServiceRepository
type ProjectValidateConsumeRepository = validateconsume.ProjectRepository
type EnvironmentValidateConsumeRepository = validateconsume.EnvironmentRepository

// ... Validate And Reservation Use Case ...

type APIKeyValidateReservationRepository = validatereservation.APIKeyRepository
type RequestValidateReservationRepository = validatereservation.RequestRepository
type ServiceValidateReservationRepository = validatereservation.ServiceRepository
type EnvironmentValidateReservationRepository = validatereservation.EnvironmentRepository
type ReservationValidateReservationRepository = validatereservation.ReservationRepository

// ... Validate Reserve Use Case ...

type RequestValidateReserveRepository = validatereserve.RequestRepository
type ReservationValidateReserveRepository = validatereserve.ReservationRepository
