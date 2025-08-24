package bootstrap

import "github.com/MAD-py/pandora-core/internal/adapters/persistence"

type Dependencies struct {
	Repositories persistence.Repositories
}

func NewDependencies(repositories persistence.Repositories) *Dependencies {
	return &Dependencies{Repositories: repositories}
}
