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
	serviceRepo := repository.NewServiceRepository(db.Pool(), db.HandlerErr())

	apiKeyUseCase := app.NewAPIKeyUseCase(apiKeyRepo, nil, nil, nil)
	serviceUseCase := app.NewServiceUseCase(serviceRepo)

	srv := http.NewServer(":8080", serviceUseCase, apiKeyUseCase)

	srv.Run()
}
