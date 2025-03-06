package main

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence/repository"
	"github.com/MAD-py/pandora-core/internal/app"
)

func main() {
	db, err := persistence.NewPersistence("host=localhost port=5436 user=postgres password=postgres dbname=pandora sslmode=disable")
	if err != nil {
		panic(err)
	}

	apiKeyRepo := repository.NewAPIKeyRepository(db.Pool(), db.HandlerErr())
	clientRepo := repository.NewClientRepository(db.Pool(), db.HandlerErr())
	serviceRepo := repository.NewServiceRepository(db.Pool(), db.HandlerErr())
	projectRepo := repository.NewProjectRepository(db.Pool(), db.HandlerErr())
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
		":8080",
		serviceUseCase,
		apiKeyUseCase,
		clientUseCase,
		projectUseCase,
		environmentUseCase,
	)

	srv.Run()
}
