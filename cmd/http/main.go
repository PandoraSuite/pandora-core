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
	log.Printf("[INFO] HTTP config loaded")

	validator := validator.NewValidator()
	log.Println("[INFO] Validator initialized")

	repositories := persistence.NewRepositories(
		persistence.PostgresDriver, cfg.DBDNS(),
	)
	log.Println("[INFO] Repositories initialized")

	jwtProvider := security.NewJWTProvider([]byte(cfg.JWTSecret()))
	log.Println("[INFO] JWT provider initialized")

	credentialsRepo := security.NewCredentialsRepository(cfg.CredentialsFile())
	log.Println("[INFO] Credentials repository initialized")

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
