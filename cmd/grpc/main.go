package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/config"
	"github.com/MAD-py/pandora-core/internal/validator"
)

func main() {
	time.Local = time.UTC

	log.Println("[INFO] Starting Pandora Core (gRPC)...")

	cfg := config.LoadGRPCConfig()
	log.Printf("[INFO] gRPC config loaded")

	validator := validator.NewValidator()
	log.Println("[INFO] Validator initialized")

	repositories := persistence.NewRepositories(
		persistence.PostgresDriver, cfg.DBDNS(),
	)
	log.Println("[INFO] Repositories initialized")

	gRPCDeps := bootstrap.NewDependencies(validator, repositories)

	srv := grpc.NewServer(
		fmt.Sprintf(":%s", cfg.Port()),
		gRPCDeps,
	)

	srv.Run()
}
