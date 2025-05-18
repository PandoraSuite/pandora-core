package commit

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id string) errors.Error
}

type useCase struct {
	validator validator.Validator

	reservationRepo ReservationRepository
}

func (uc *useCase) Execute(ctx context.Context, id string) errors.Error {
	if err := uc.validateID(id); err != nil {
		return err
	}
	return uc.reservationRepo.Delete(ctx, id)
}

func (uc *useCase) validateID(id string) errors.Error {
	return uc.validator.ValidateVariable(
		id,
		"id",
		"required,uuid4",
		map[string]string{
			"uuid4":    "id must be a valid UUID",
			"required": "id is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator,
	reservationRepo ReservationRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		reservationRepo: reservationRepo,
	}
}
