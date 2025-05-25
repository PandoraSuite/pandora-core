package service

import (
	"github.com/MAD-py/pandora-core/internal/app/service/create"
	"github.com/MAD-py/pandora-core/internal/app/service/delete"
	"github.com/MAD-py/pandora-core/internal/app/service/list"
	listrequest "github.com/MAD-py/pandora-core/internal/app/service/list_request"
	updatestatus "github.com/MAD-py/pandora-core/internal/app/service/update_status"
)

// ... Create Use Case ...

type ServiceCreateRepository = create.ServiceRepository

// ... Delete Use Case ...

type ServiceDeleteRepository = delete.ServiceRepository
type ProjectServiceVerifier = delete.ProjectRepository

// ... List Use Case ...

type ServiceListRepository = list.ServiceRepository

// ... List Requests Use Case ...
type ServiceListRequestsRepository = listrequest.ServiceRepository
type RequestListByServiceRepository = listrequest.RequestRepository

// ... Update Status Use Case ...

type ServiceUpdateStatusRepository = updatestatus.ServiceRepository
