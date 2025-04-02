package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc/api_key"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

type Server struct {
	addr string

	server *grpc.Server

	apiKeyService inbound.APIKeyGRPCPort
}

func (s *Server) setupServices() {
	api_key.RegisterService(s.server, s.apiKeyService)
}

func (s *Server) Run() {
	s.server = grpc.NewServer()
	s.setupServices()

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("[ERROR] Failed to listen: %v", err)
	}

	log.Printf("[INFO] gRPC server is running on port: %s\n", s.addr)
	log.Printf("[INFO] Pandora Core is fully initialized and ready to accept requests.\n\n")
	if err := s.server.Serve(listener); err != nil {
		log.Fatalf("[ERROR] Failed to serve: %v", err)
	}
}

func NewServer(
	addr string,
	apiKeyService inbound.APIKeyGRPCPort,
) *Server {
	return &Server{
		addr:          addr,
		apiKeyService: apiKeyService,
	}
}
