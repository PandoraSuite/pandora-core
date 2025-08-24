package tasks

import (
	"github.com/MAD-py/go-taskengine/taskengine"
	"github.com/MAD-py/pandora-core/internal/adapters/taskengine/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/taskengine/jobs"
	"github.com/MAD-py/pandora-core/internal/app/project"
)

func ProjectQuotaReset(deps *bootstrap.Dependencies) (*taskengine.Task, error) {
	resetDueRequestsUseCase := project.NewResetDueRequestsUseCase(deps.Repositories.Project())
	return taskengine.NewTask(
		"project-quota-reset",
		jobs.ProjectQuotaReset(resetDueRequestsUseCase),
	)
}
