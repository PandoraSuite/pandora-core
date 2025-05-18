package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/persistence"
	"github.com/MAD-py/pandora-core/internal/validator"
)

func main() {
	time.Local = time.UTC

	log.Println("[INFO] Starting Pandora Core (gRPC)...")

	validator := validator.NewValidator()

	repositories := persistence.NewRepositories(
		persistence.PostgresDriver,
		"postgresql://postgres:postgres@localhost:5436/pandora",
	)

	gRPCDeps := bootstrap.NewDependencies(validator, repositories)

	srv := grpc.NewServer(
		fmt.Sprintf(":%s", "50051"),
		gRPCDeps,
	)

	srv.Run()
}
