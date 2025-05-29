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
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type service struct {
	pb.APIKeyServiceServer

	validateUC        apikey.ValidateUseCase
	validateConsumeUC apikey.ValidateConsumeUseCase
}

func (s *service) Validate(
	ctx context.Context, req *pb.ValidateRequest,
) (*pb.ValidateResponse, error) {
	reqValidate := dto.APIKeyValidate{
		APIKey: req.ApiKey,
		Request: &dto.RequestCreate{
			Path:   req.Request.Path,
			Method: req.Request.Method,
			Metadata: &dto.RequestMetadata{
				Body:            req.Request.Metadata.Body,
				Headers:         req.Request.Metadata.Headers,
				QueryParams:     req.Request.Metadata.QueryParams,
				BodyContentType: enums.RequestBodyContentType(req.Request.Metadata.BodyContentType),
			},
			IPAddress:   req.Request.IpAddress,
			RequestTime: req.Request.RequestTime.AsTime(),
		},
		ServiceName:    req.ServiceName,
		ServiceVersion: req.ServiceVersion,
	}

	response, err := s.validateUC.Execute(ctx, &reqValidate)
	if err != nil {
		return nil, status.Error(
			utils.GetDomainErrorStatusCode(err),
			err.Error(),
		)
	}

	return &pb.ValidateResponse{
		Valid:       response.Valid,
		RequestId:   response.RequestID,
		FailureCode: string(response.FailureCode),
		ConsumerInfo: &pb.ConsumerInfo{
			ProjectId:   int64(response.ConsumerInfo.ProjectID),
			ProjectName: response.ConsumerInfo.ProjectName,
			ClientId:    int64(response.ConsumerInfo.ClientID),
			ClientName:  response.ConsumerInfo.ClientName,
		},
	}, nil
}

func (s *service) ValidateConsume(
	ctx context.Context, req *pb.ValidateRequest,
) (*pb.ValidateConsumeResponse, error) {
	reqValidate := dto.APIKeyValidate{
		APIKey: req.ApiKey,
		Request: &dto.RequestCreate{
			Path:   req.Request.Path,
			Method: req.Request.Method,
			Metadata: &dto.RequestMetadata{
				Body:            req.Request.Metadata.Body,
				Headers:         req.Request.Metadata.Headers,
				QueryParams:     req.Request.Metadata.QueryParams,
				BodyContentType: enums.RequestBodyContentType(req.Request.Metadata.BodyContentType),
			},
			IPAddress:   req.Request.IpAddress,
			RequestTime: req.Request.RequestTime.AsTime(),
		},
		ServiceName:    req.ServiceName,
		ServiceVersion: req.ServiceVersion,
	}

	response, err := s.validateConsumeUC.Execute(ctx, &reqValidate)
	if err != nil {
		return nil, status.Error(
			utils.GetDomainErrorStatusCode(err),
			err.Error(),
		)
	}

	return &pb.ValidateConsumeResponse{
		BaseResponse: &pb.ValidateResponse{
			Valid:       response.Valid,
			RequestId:   response.RequestID,
			FailureCode: string(response.FailureCode),
			ConsumerInfo: &pb.ConsumerInfo{
				ProjectId:   int64(response.ConsumerInfo.ProjectID),
				ProjectName: response.ConsumerInfo.ProjectName,
				ClientId:    int64(response.ConsumerInfo.ClientID),
				ClientName:  response.ConsumerInfo.ClientName,
			},
		},
		AvailableRequest: int64(response.AvailableRequest),
	}, nil
}

func RegisterService(s *grpc.Server, deps *bootstrap.Dependencies) {
	service := service{
		validateUC: apikey.NewValidateUseCase(
			deps.Validator,
			deps.Repositories.APIKey(),
			deps.Repositories.Project(),
			deps.Repositories.Service(),
			deps.Repositories.Request(),
			deps.Repositories.Environment(),
		),
		validateConsumeUC: apikey.NewValidateConsumeUseCase(
			deps.Validator,
			deps.Repositories.APIKey(),
			deps.Repositories.Project(),
			deps.Repositories.Request(),
			deps.Repositories.Service(),
			deps.Repositories.Environment(),
		),
	}
	pb.RegisterAPIKeyServiceServer(s, &service)
}
