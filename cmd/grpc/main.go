package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence/repository"
	"github.com/MAD-py/pandora-core/internal/app"
)

func main() {
	time.Local = time.UTC

	log.Println("[INFO] Starting Pandora Core (gRPC)...")
	db, err := persistence.NewPersistence("postgresql://postgres:postgres@localhost:5436/pandora")
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize persistence: %v", err)
	}
	pool := db.Pool()
	errHandler := db.HandlerErr()

	apiKeyRepo := repository.NewAPIKeyRepository(pool, errHandler)
	serviceRepo := repository.NewServiceRepository(pool, errHandler)
	requestLogRepo := repository.NewRequestLogRepository(pool, errHandler)
	environmentRepo := repository.NewEnvironmentRepository(pool, errHandler)
	reservationRepo := repository.NewReservationRepository(pool, errHandler)

	apiKeyUseCase := app.NewAPIKeyUseCase(
		apiKeyRepo, requestLogRepo, serviceRepo, environmentRepo, reservationRepo,
	)
	requestLogUseCase := app.NewRequestLogUseCase(requestLogRepo)
	reservationUseCase := app.NewReservationUseCase(reservationRepo, environmentRepo)

	srv := grpc.NewServer(
		fmt.Sprintf(":%s", "50051"),
		apiKeyUseCase,
		requestLogUseCase,
		reservationUseCase,
	)

	srv.Run()
}
