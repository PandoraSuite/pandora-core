package main

import (
	"fmt"

	"github.com/MAD-py/pandora-core/cmd/http/config"
	"github.com/MAD-py/pandora-core/internal/adapters/http"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence/repository"
	"github.com/MAD-py/pandora-core/internal/adapters/security"
	"github.com/MAD-py/pandora-core/internal/app"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := persistence.NewPersistence(config.DBDNS())
	if err != nil {
		panic(err)
	}

	apiKeyRepo := repository.NewAPIKeyRepository(db.Pool(), db.HandlerErr())
	clientRepo := repository.NewClientRepository(db.Pool(), db.HandlerErr())
	serviceRepo := repository.NewServiceRepository(db.Pool(), db.HandlerErr())
	projectRepo := repository.NewProjectRepository(db.Pool(), db.HandlerErr())
	jwtProvider := security.NewJWTProvider([]byte(config.JWTSecret()))
	requestLogRepo := repository.NewRequestLogRepository(
		db.Pool(), db.HandlerErr(),
	)
	environmentRepo := repository.NewEnvironmentRepository(
		db.Pool(), db.HandlerErr(),
	)
	projectServiceRepo := repository.NewProjectServiceRepository(
		db.Pool(), db.HandlerErr(),
	)
	environmentServiceRepo := repository.NewEnvironmentServiceRepository(
		db.Pool(), db.HandlerErr(),
	)

	credentialsFile, err := config.CredentialsFile()
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
	projectUseCase := app.NewProjectUseCase(projectRepo, projectServiceRepo)
	environmentUseCase := app.NewEnvironmentUseCase(
		environmentRepo, projectServiceRepo, environmentServiceRepo,
	)
	apiKeyUseCase := app.NewAPIKeyUseCase(
		apiKeyRepo, requestLogRepo, serviceRepo, environmentServiceRepo,
	)

	srv := http.NewServer(
		fmt.Sprintf(":%s", config.HTTPPort()),
		serviceUseCase,
		authUseCase,
		apiKeyUseCase,
		clientUseCase,
		projectUseCase,
		environmentUseCase,
	)

	srv.Run()
}
