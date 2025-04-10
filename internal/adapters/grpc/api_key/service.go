package api_key

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	pb "github.com/MAD-py/pandora-core/internal/adapters/grpc/api_key/v1"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

type service struct {
	pb.APIKeyServiceServer

	apiKeyService inbound.APIKeyGRPCPort
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

func (s *service) ValidateAndReserve(ctx context.Context, req *pb.ValidateAndReserveRequest) (*pb.ValidateAndReserveResponse, error) {
	params := req.GetParams()
	reqValidate := dto.APIKeyValidate{
		Key:            params.Key,
		Service:        params.Service,
		Environment:    params.Environment,
		ServiceVersion: params.ServiceVersion,
		RequestTime:    params.RequestTime.AsTime(),
	}
	response, err := s.apiKeyService.ValidateAndReserve(ctx, &reqValidate)
	if err != nil {
		return nil, status.Error(
			utils.GetDomainErrorStatusCode(err),
			err.Message,
		)
	}
	if response.Valid {
		return &pb.ValidateAndReserveResponse{
			Valid: true,
			Result: &pb.ValidateAndReserveResponse_Successful_{
				Successful: &pb.ValidateAndReserveResponse_Successful{
					RequestId:        response.RequestID,
					ReservationId:    response.ReservationID,
					AvailableRequest: int64(response.AvailableRequest),
				},
			},
		}, nil
	} else {
		return &pb.ValidateAndReserveResponse{
			Valid: false,
			Result: &pb.ValidateAndReserveResponse_Failed_{
				Failed: &pb.ValidateAndReserveResponse_Failed{
					Code:    response.Code.String(),
					Message: response.Message,
				},
			},
		}, nil
	}
}

func (s *service) ValidateWithReservationRequest(ctx context.Context, req *pb.ValidateWithReservationRequest) (*pb.ValidateWithReservationResponse, error) {
	return &pb.ValidateWithReservationResponse{
		Valid: true,
		Result: &pb.ValidateWithReservationResponse_Successful_{
			Successful: &pb.ValidateWithReservationResponse_Successful{
				RequestId: "req-123456",
			},
		},
	}, nil
}

func RegisterService(server *grpc.Server, apiKeyService inbound.APIKeyGRPCPort) {
	service := &service{apiKeyService: apiKeyService}
	pb.RegisterAPIKeyServiceServer(server, service)
}
