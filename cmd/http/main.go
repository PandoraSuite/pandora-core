package main

import (
	"fmt"
	"log"

	"github.com/MAD-py/pandora-core/cmd/http/config"
	"github.com/MAD-py/pandora-core/internal/adapters/http"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence/repository"
	"github.com/MAD-py/pandora-core/internal/adapters/security"
	"github.com/MAD-py/pandora-core/internal/app"
)

func main() {
	log.Println("[INFO] Starting Pandora Core...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[ERROR] Error loading configuration: %v", err)
	}

	log.Println("[INFO] Configuration loaded successfully from environment variables and configuration files.")

	db, err := persistence.NewPersistence(cfg.DBDNS())
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize persistence: %v", err)
	}
	pool := db.Pool()
	errHandler := db.HandlerErr()

	apiKeyRepo := repository.NewAPIKeyRepository(pool, errHandler)
	clientRepo := repository.NewClientRepository(pool, errHandler)
	serviceRepo := repository.NewServiceRepository(pool, errHandler)
	projectRepo := repository.NewProjectRepository(pool, errHandler)
	requestLogRepo := repository.NewRequestLogRepository(pool, errHandler)
	environmentRepo := repository.NewEnvironmentRepository(pool, errHandler)

	jwtProvider := security.NewJWTProvider([]byte(cfg.JWTSecret()))
	credentialsFile, err := cfg.CredentialsFile()
	if err != nil {
		panic(err)
	}

	credentialsRepo, err := security.NewCredentialsRepository(credentialsFile)
	if err != nil {
		panic(err)
	}

	authUseCase := app.NewAuthUseCase(jwtProvider, credentialsRepo)
	clientUseCase := app.NewClientUseCase(clientRepo)
	serviceUseCase := app.NewServiceUseCase(serviceRepo)
	projectUseCase := app.NewProjectUseCase(projectRepo)
	environmentUseCase := app.NewEnvironmentUseCase(environmentRepo)
	apiKeyUseCase := app.NewAPIKeyUseCase(
		apiKeyRepo, requestLogRepo, serviceRepo, environmentRepo,
	)

	srv := http.NewServer(
		fmt.Sprintf(":%s", cfg.HTTPPort()),
		serviceUseCase,
		authUseCase,
		apiKeyUseCase,
		clientUseCase,
		projectUseCase,
		environmentUseCase,
	)

	srv.Run()
}
