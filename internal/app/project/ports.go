package project

import (
	assignservice "github.com/MAD-py/pandora-core/internal/app/project/assign_service"
	"github.com/MAD-py/pandora-core/internal/app/project/create"
	"github.com/MAD-py/pandora-core/internal/app/project/delete"
	"github.com/MAD-py/pandora-core/internal/app/project/get"
	"github.com/MAD-py/pandora-core/internal/app/project/list"
	listenvironments "github.com/MAD-py/pandora-core/internal/app/project/list_environments"
	removeservice "github.com/MAD-py/pandora-core/internal/app/project/remove_service"
	resetduerequests "github.com/MAD-py/pandora-core/internal/app/project/reset_due_requests"
	resetrequests "github.com/MAD-py/pandora-core/internal/app/project/reset_requests"
	"github.com/MAD-py/pandora-core/internal/app/project/update"
	updateservice "github.com/MAD-py/pandora-core/internal/app/project/update_service"
)

// ... Assign Service Use Case ...

type ProjectAssignServiceRepository = assignservice.ProjectRepository

// ... Create Use Case ...

type ProjectCreateRepository = create.ProjectRepository

// ... Delete Use Case ...

type ProjectDeleteRepository = delete.ProjectRepository

// ... Get Use Case ...

type ProjectGetRepository = get.ProjectRepository

// ... List Use Case ...

type ProjectListRepository = list.ProjectRepository

// ... List Environments Use Case ...

type ProjectListEnvironmentsRepository = listenvironments.ProjectRepository
type EnvironmentListByClientRepository = listenvironments.EnvironmentRepository

// ... Remove Service Use Case ...

type ProjectRemoveServiceRepository = removeservice.ProjectRepository
type EnvironmentRemoveServiceRepository = removeservice.EnvironmentRepository

// ... Reset Request Use Case ...

type ProjectResetRequestRepository = resetrequests.ProjectRepository

// ... Reset Due Requests Use Case ...

type ProjectResetDueRequestsRepository = resetduerequests.ProjectRepository

// ... Update Use Case ...

type ProjectUpdateRepository = update.ProjectRepository

// ... Update Service Use Case ...

type ProjectUpdateServiceRepository = updateservice.ProjectRepository
type EnvironmentServiceInfiniteQuotaRepository = updateservice.EnvironmentRepository
