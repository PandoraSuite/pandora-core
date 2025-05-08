package bootstrap

import (
	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/adapters/security"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type Dependencies struct {
	Validator validator.Validator

	TokenProvider security.TokenProvider

	Repositories    persistence.Repositories
	CredentialsRepo security.CredentialsRepository
}

func NewDependencies(
	validator validator.Validator,
	repositories persistence.Repositories,
	tokenProvider security.TokenProvider,
	credentialsRepo security.CredentialsRepository,
) *Dependencies {
	return &Dependencies{
		Validator:       validator,
		Repositories:    repositories,
		TokenProvider:   tokenProvider,
		CredentialsRepo: credentialsRepo,
	}
}
