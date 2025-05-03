package reservation

import (
	"github.com/MAD-py/pandora-core/internal/app/reservation/commit"
	"github.com/MAD-py/pandora-core/internal/app/reservation/rollback"
	"github.com/MAD-py/pandora-core/internal/validator"
)

// ... Commit Use Case ...

type CommitUseCase = commit.UseCase

func NewCommitUseCase(
	validator validator.Validator,
	reservationRepo ReservationCommitRepository,
) CommitUseCase {
	return commit.NewUseCase(validator, reservationRepo)
}

// ... Rollback Use Case ...

type RollbackUseCase = rollback.UseCase

func NewRollbackUseCase(
	validator validator.Validator,
	reservationRepo ReservationRollbackRepository,
	environmentRepo EnvironmentAvailableRequestIncrementerRepository,
) RollbackUseCase {
	return rollback.NewUseCase(validator, reservationRepo, environmentRepo)
}
