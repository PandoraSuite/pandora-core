package apikey

import (
	"github.com/MAD-py/pandora-core/internal/app/api_key/create"
	"github.com/MAD-py/pandora-core/internal/app/api_key/delete"
	revealkey "github.com/MAD-py/pandora-core/internal/app/api_key/reveal_key"
	"github.com/MAD-py/pandora-core/internal/app/api_key/update"
	validateconsume "github.com/MAD-py/pandora-core/internal/app/api_key/validate_consume"
	validateonly "github.com/MAD-py/pandora-core/internal/app/api_key/validate_only"
)

// ... Create Use Case ...

type APIKeyCreateRepository = create.APIKeyRepository

// ... Delete Use Case ...

type APIKeyDeleteRepository = delete.APIKeyRepository

// ... Update Use Case ...

type APIKeyUpdateRepository = update.APIKeyRepository

// ... Validate Use Case ...

type APIKeyValidateRepository = validateonly.APIKeyRepository
type ProjectValidateRepository = validateonly.ProjectRepository
type ServiceValidateRepository = validateonly.ServiceRepository
type RequestValidateRepository = validateonly.RequestRepository
type EnvironmentValidateRepository = validateonly.EnvironmentRepository

// ... Validate And Consume Use Case ...

type APIKeyValidateConsumeRepository = validateconsume.APIKeyRepository
type RequestValidateConsumeRepository = validateconsume.RequestRepository
type ServiceValidateConsumeRepository = validateconsume.ServiceRepository
type ProjectValidateConsumeRepository = validateconsume.ProjectRepository
type EnvironmentValidateConsumeRepository = validateconsume.EnvironmentRepository

// ... Reveal Key Use Case ...

type APIKeyRevealKeyRepository = revealkey.APIKeyRepository
