package bootstrap

import (
	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/ports"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type Dependencies struct {
	Validator validator.Validator

	TokenProvider ports.TokenProvider

	Repositories    persistence.Repositories
	CredentialsRepo ports.CredentialsRepository
}

func NewDependencies(
	validator validator.Validator,
	repositories persistence.Repositories,
	tokenProvider ports.TokenProvider,
	credentialsRepo ports.CredentialsRepository,
) *Dependencies {
	return &Dependencies{
		Validator:       validator,
		Repositories:    repositories,
		TokenProvider:   tokenProvider,
		CredentialsRepo: credentialsRepo,
	}
}
