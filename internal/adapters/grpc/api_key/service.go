package api_key

import (
	"context"

	"google.golang.org/grpc"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc/api_key/pb"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

type service struct {
	pb.APIKeyServiceServer

	apiKeyService inbound.APIKeyGRPCPort
}

func (s *service) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	return &pb.ValidateResponse{
		Valid: true,
		Result: &pb.ValidateResponse_Successful_{
			Successful: &pb.ValidateResponse_Successful{
				RequestId: "req-123456",
			},
		},
	}, nil
}

func (s *service) ValidateAndConsume(ctx context.Context, req *pb.ValidateAndConsumeRequest) (*pb.ValidateAndConsumeResponse, error) {
	return &pb.ValidateAndConsumeResponse{
		Valid: true,
		Result: &pb.ValidateAndConsumeResponse_Successful_{
			Successful: &pb.ValidateAndConsumeResponse_Successful{
				RequestId:        "req-123456",
				AvailableRequest: 1000,
			},
		},
	}, nil
}

func (s *service) ValidateAndBooking(ctx context.Context, req *pb.ValidateAndReservationRequest) (*pb.ValidateAndReservationResponse, error) {
	return &pb.ValidateAndReservationResponse{
		Valid: true,
		Result: &pb.ValidateAndReservationResponse_Successful_{
			Successful: &pb.ValidateAndReservationResponse_Successful{
				RequestId:        "req-123456",
				ReservationId:    "booking-123456",
				AvailableRequest: 1000,
			},
		},
	}, nil
}

func RegisterService(server *grpc.Server, apiKeyService inbound.APIKeyGRPCPort) {
	service := &service{apiKeyService: apiKeyService}
	pb.RegisterAPIKeyServiceServer(server, service)
}
