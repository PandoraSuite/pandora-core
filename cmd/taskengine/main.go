package main

import (
	"log"
	"time"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/adapters/taskengine"
	"github.com/MAD-py/pandora-core/internal/adapters/taskengine/bootstrap"
	"github.com/MAD-py/pandora-core/internal/config"
)

func main() {
	time.Local = time.UTC

	log.Println("[INFO] Starting Pandora Core (TaskEngine)...")

	cfg := config.LoadTaskEngineConfig()
	log.Printf("[INFO] TaskEngine config loaded")

	repositories := persistence.NewRepositories(
		persistence.PostgresDriver, cfg.DBDNS(),
	)
	log.Println("[INFO] Repositories initialized")

	taskEngineDeps := bootstrap.NewDependencies(repositories)
	log.Println("[INFO] TaskEngine dependencies initialized")

	engine, err := taskengine.NewEngine(cfg.DBDNS(), taskEngineDeps)
	if err != nil {
		log.Fatalf("[ERROR] Failed to create TaskEngine: %v", err)
	}
	log.Println("[INFO] TaskEngine created successfully")

	engine.Run()
}
