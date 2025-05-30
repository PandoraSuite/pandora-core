package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MAD-py/pandora-core/cmd/grpc/config"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence/repository"
	"github.com/MAD-py/pandora-core/internal/app"
)

func main() {
	time.Local = time.UTC

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[ERROR] Error loading configuration: %v", err)
	}

	log.Println("[INFO] Starting Pandora Core (gRPC)...")
	db, err := persistence.NewPersistence(cfg.DBDNS())
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize persistence: %v", err)
	}
	pool := db.Pool()
	errHandler := db.HandlerErr()

	apiKeyRepo := repository.NewAPIKeyRepository(pool, errHandler)
	serviceRepo := repository.NewServiceRepository(pool, errHandler)
	projectRepo := repository.NewProjectRepository(pool, errHandler)
	requestLogRepo := repository.NewRequestLogRepository(pool, errHandler)
	environmentRepo := repository.NewEnvironmentRepository(pool, errHandler)
	reservationRepo := repository.NewReservationRepository(pool, errHandler)

	apiKeyUseCase := app.NewAPIKeyUseCase(
		apiKeyRepo, requestLogRepo, serviceRepo, projectRepo, environmentRepo, reservationRepo,
	)
	requestLogUseCase := app.NewRequestLogUseCase(requestLogRepo)
	reservationUseCase := app.NewReservationUseCase(reservationRepo, environmentRepo)

	srv := grpc.NewServer(
		fmt.Sprintf(":%s", cfg.Port()),
		apiKeyUseCase,
		requestLogUseCase,
		reservationUseCase,
	)

	srv.Run()
}
