package bootstrap

import (
	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type Dependencies struct {
	Validator validator.Validator

	Repositories persistence.Repositories
}

func NewDependencies(
	validator validator.Validator, repositories persistence.Repositories,
) *Dependencies {
	return &Dependencies{
		Validator:    validator,
		Repositories: repositories,
	}
}
