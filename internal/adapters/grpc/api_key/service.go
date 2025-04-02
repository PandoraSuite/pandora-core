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

func (s *service) Validate(ctx context.Context, req *pb.APIKeyValidateBooking) (*pb.APIKeyValidateResponse, error) {
	return &pb.APIKeyValidateResponse{
		Valid: true,
		Result: &pb.APIKeyValidateResponse_Successful_{
			Successful: &pb.APIKeyValidateResponse_Successful{
				RequestId: "req-123456",
			},
		},
	}, nil
}

func (s *service) ValidateAndConsume(ctx context.Context, req *pb.APIKeyValidate) (*pb.APIKeyValidateConsumeResponse, error) {
	return &pb.APIKeyValidateConsumeResponse{
		Valid: true,
		Result: &pb.APIKeyValidateConsumeResponse_Successful_{
			Successful: &pb.APIKeyValidateConsumeResponse_Successful{
				RequestId:        "req-123456",
				AvailableRequest: 1000,
			},
		},
	}, nil
}

func (s *service) ValidateAndBooking(ctx context.Context, req *pb.APIKeyValidate) (*pb.APIKeyValidateBookingResponse, error) {
	return &pb.APIKeyValidateBookingResponse{
		Valid: true,
		Result: &pb.APIKeyValidateBookingResponse_Successful_{
			Successful: &pb.APIKeyValidateBookingResponse_Successful{
				RequestId:        "req-123456",
				BookingId:        "booking-123456",
				AvailableRequest: 1000,
			},
		},
	}, nil
}

func RegisterService(server *grpc.Server, apiKeyService inbound.APIKeyGRPCPort) {
	service := &service{apiKeyService: apiKeyService}
	pb.RegisterAPIKeyServiceServer(server, service)
}
