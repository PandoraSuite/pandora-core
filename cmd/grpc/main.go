package main

import (
	"fmt"
	"log"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc"
)

func main() {
	log.Println("[INFO] Starting Pandora Core (gRPC)...")

	srv := grpc.NewServer(fmt.Sprintf(":%s", "50051"))
	srv.Run()
}
