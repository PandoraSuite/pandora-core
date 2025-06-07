package environment

import (
	assignservice "github.com/MAD-py/pandora-core/internal/app/environment/assign_service"
	"github.com/MAD-py/pandora-core/internal/app/environment/create"
	"github.com/MAD-py/pandora-core/internal/app/environment/delete"
	"github.com/MAD-py/pandora-core/internal/app/environment/get"
	listapikey "github.com/MAD-py/pandora-core/internal/app/environment/list_api_key"
	removeservice "github.com/MAD-py/pandora-core/internal/app/environment/remove_service"
	resetrequests "github.com/MAD-py/pandora-core/internal/app/environment/reset_requests"
	"github.com/MAD-py/pandora-core/internal/app/environment/update"
	updateservice "github.com/MAD-py/pandora-core/internal/app/environment/update_service"
)

// ... Assign Service Use Case ...

type EnvironmentAssignServiceRepository = assignservice.EnvironmentRepository

// ... Create Use Case ...

type EnvironmentCreateRepository = create.EnvironmentRepository
type ProjectQuotaRepository = create.ProjectRepository

// ... Delete Use Case ...
type EnvironmentDeleteRepository = delete.EnvironmentRepository

// ... Get Use Case ...

type EnvironmentGetRepository = get.EnvironmentRepository

// ... List API Key Use Case ...

type EnvironmentListAPIKeyRepository = listapikey.EnvironmentRepository
type APIKeyListByEnvironmentRepository = listapikey.APIKeyRepository

// ... Remove Service Use Case ...

type EnvironmentRemoveServiceRepository = removeservice.EnvironmentRepository

// ... Reset Request Use Case ...

type EnvironmentResetRequestRepository = resetrequests.EnvironmentRepository

// ... Update Use Case ...

type EnvironmentUpdateRepository = update.EnvironmentRepository

// ... Update Service Use Case ...

type EnvironmentUpdateServiceRepository = updateservice.EnvironmentRepository
