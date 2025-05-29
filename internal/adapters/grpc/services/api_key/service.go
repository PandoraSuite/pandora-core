package apikey

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc/bootstrap"
	pb "github.com/MAD-py/pandora-core/internal/adapters/grpc/services/api_key/v1"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/utils"
	apikey "github.com/MAD-py/pandora-core/internal/app/api_key"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type service struct {
	pb.APIKeyServiceServer

	validateConsumeUC apikey.ValidateConsumeUseCase
}

func (s *service) ValidateAndConsume(ctx context.Context, req *pb.ValidateAndConsumeRequest) (*pb.ValidateAndConsumeResponse, error) {
	params := req.GetParams()
	reqValidate := dto.APIKeyValidate{
		APIKey:         params.Key,
		Service:        params.Service,
		Environment:    params.Environment,
		ServiceVersion: params.ServiceVersion,
		RequestTime:    params.RequestTime.AsTime(),
	}
	response, err := s.validateAndConsumeUC.Execute(ctx, &reqValidate)
	if err != nil {
		return nil, status.Error(
			utils.GetDomainErrorStatusCode(err),
			err.Error(),
		)
	}
	if response.Valid {
		return &pb.ValidateAndConsumeResponse{
			Valid: true,
			Result: &pb.ValidateAndConsumeResponse_Successful_{
				Successful: &pb.ValidateAndConsumeResponse_Successful{
					RequestId:        response.RequestID,
					AvailableRequest: int64(response.AvailableRequest),
				},
			},
		}, nil
	} else {
		return &pb.ValidateAndConsumeResponse{
			Valid: false,
			Result: &pb.ValidateAndConsumeResponse_Failed_{
				Failed: &pb.ValidateAndConsumeResponse_Failed{
					Code:    response.Code.String(),
					Message: response.Message,
				},
			},
		}, nil
	}
}

func RegisterService(s *grpc.Server, deps *bootstrap.Dependencies) {
	service := service{

		validateAndConsumeUC: apikey.NewValidateConsumeUseCase(
			deps.Validator,
			deps.Repositories.APIKey(),
			deps.Repositories.Request(),
			deps.Repositories.Service(),
			deps.Repositories.Reservation(),
			deps.Repositories.Environment(),
		),
	}
	pb.RegisterAPIKeyServiceServer(s, &service)
}
