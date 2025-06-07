package client

import (
	"github.com/MAD-py/pandora-core/internal/app/client/create"
	"github.com/MAD-py/pandora-core/internal/app/client/get"
	"github.com/MAD-py/pandora-core/internal/app/client/list"
	listprojects "github.com/MAD-py/pandora-core/internal/app/client/list_projects"
	"github.com/MAD-py/pandora-core/internal/app/client/update"
)

// ... Create Use Case ...

type ClientCreateRepository = create.ClientRepository

// ... Delete Use Case ...

type ClientDeleteRepository = create.ClientRepository

// ... Get Use Case ...

type ClientGetRepository = get.ClientRepository

// ... List Use Case ...

type ClientListRepository = list.ClientRepository

// ... List Projects Use Case ...

type ClientListProjectsRepository = listprojects.ClientRepository
type ProjectListByClientRepository = listprojects.ProjectRepository

// ... Update Use Case ...
type ClientUpdateRepository = update.ClientRepository
