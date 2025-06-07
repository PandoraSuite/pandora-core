package delete

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int) errors.Error
}

type useCase struct {
	validator validator.Validator

	clientRepo ClientRepository
}

func (u *useCase) Execute(ctx context.Context, id int) errors.Error {
	if err := u.validateID(id); err != nil {
		return err
	}

	if err := u.clientRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (uc *useCase) validateID(id int) errors.Error {
	return uc.validator.ValidateVariable(
		id,
		"id",
		"required,gt=0",
		map[string]string{
			"gt":       "id must be greater than 0",
			"required": "id is required",
		},
	)
}
func NewUseCase(
	validator validator.Validator,
	clientRepo ClientRepository,
) UseCase {
	return &useCase{
		validator:  validator,
		clientRepo: clientRepo,
	}
}
