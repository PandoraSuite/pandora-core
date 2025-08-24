package taskengine

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/MAD-py/go-taskengine/taskengine"
	"github.com/MAD-py/go-taskengine/taskengine/store/postgresql"
	"github.com/MAD-py/pandora-core/internal/adapters/taskengine/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/taskengine/registry"
	"github.com/MAD-py/pandora-core/internal/adapters/taskengine/tasks"
)

type Engine struct {
	engine *taskengine.Engine
	deps   *bootstrap.Dependencies
}

func (e *Engine) Run() error {
	{
		task, err := tasks.ProjectQuotaReset(e.deps)
		if err != nil {
			log.Printf("[ERROR] Failed to create project quota reset task: %v\n", err)
			return err
		}

		err = registry.ProjectQuotaReset(e.engine, task)
		if err != nil {
			log.Printf("[ERROR] Failed to register project quota reset task: %v\n", err)
			return err
		}
	}

	log.Printf("[INFO] Task Engine is starting...")
	if err := e.engine.Run(); err != nil {
		log.Printf("[ERROR] Failed to start server: %v\n", err)
		return err
	}

	return nil
}

func NewEngine(connString string, deps *bootstrap.Dependencies) (*Engine, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to task engine database: %v\n", err)
		return nil, err
	}

	engine, err := taskengine.New(postgresql.NewStore(db))
	if err != nil {
		return nil, err
	}

	return &Engine{
		engine: engine,
		deps:   deps,
	}, nil
}
