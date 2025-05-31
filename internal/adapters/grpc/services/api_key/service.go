package apikey

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/MAD-py/pandora-core/internal/adapters/grpc/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/grpc/errors"
	pb "github.com/MAD-py/pandora-core/internal/adapters/grpc/services/api_key/v1"
	apikey "github.com/MAD-py/pandora-core/internal/app/api_key"
)

type service struct {
	pb.APIKeyServiceServer

	validateUC        apikey.ValidateUseCase
	validateConsumeUC apikey.ValidateConsumeUseCase
}

func (s *service) Validate(
	ctx context.Context, req *pb.ValidateRequest,
) (*pb.ValidateResponse, error) {
	response, err := s.validateUC.Execute(ctx, validateRequestToDomain(req))
	if err != nil {
		return nil, status.Error(
			errors.CodeToGRPCCode(err.Code()),
			err.Error(),
		)
	}
	return validateResponseFromDomain(response), nil
}

func (s *service) ValidateConsume(
	ctx context.Context, req *pb.ValidateRequest,
) (*pb.ValidateConsumeResponse, error) {
	response, err := s.validateConsumeUC.Execute(ctx, validateRequestToDomain(req))
	if err != nil {
		return nil, status.Error(
			errors.CodeToGRPCCode(err.Code()),
			err.Error(),
		)
	}
	return validateConsumeResponseFromDomain(response), nil
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
