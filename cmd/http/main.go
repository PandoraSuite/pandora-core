package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MAD-py/pandora-core/internal/adapters/http"
	"github.com/MAD-py/pandora-core/internal/adapters/http/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/adapters/security"
	"github.com/MAD-py/pandora-core/internal/config"
	"github.com/MAD-py/pandora-core/internal/validator"
)

func main() {
	time.Local = time.UTC

	log.Println("[INFO] Starting Pandora Core (API RESTful)...")

	cfg := config.LoadHTTPConfig()
	// if err != nil {
	// 	log.Fatalf("[ERROR] Error loading configuration: %v", err)
	// }

	log.Println("[INFO] Configuration loaded successfully from environment variables and configuration files.")

	validator := validator.NewValidator()

	repositories := persistence.NewRepositories(
		persistence.PostgresDriver, cfg.DBDNS(),
	)
	// if err != nil {
	// 	log.Fatalf("[ERROR] Failed to initialize persistence: %v", err)
	// }

	jwtProvider := security.NewJWTProvider([]byte(cfg.JWTSecret()))

	credentialsFile := cfg.CredentialsFile()
	// if err != nil {
	// 	log.Fatalf("[ERROR] Failed to load credentials file: %v", err)
	// }

	credentialsRepo := security.NewCredentialsRepository(credentialsFile)
	// if err != nil {
	// 	log.Fatalf("[ERROR] Failed to initialize credentials repository: %v", err)
	// }

	httpDeps := bootstrap.NewDependencies(
		validator,
		repositories,
		jwtProvider,
		credentialsRepo,
	)

	srv := http.NewServer(
		fmt.Sprintf(":%s", cfg.Port()),
		cfg.ExposeVersion(),
		httpDeps,
	)

	srv.Run()
}
