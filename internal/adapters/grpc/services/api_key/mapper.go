package apikey

import (
	pb "github.com/MAD-py/pandora-core/internal/adapters/grpc/services/api_key/v1"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

func validateRequestToDomain(req *pb.ValidateRequest) *dto.APIKeyValidate {
	var request *dto.RequestIncoming
	if req.Request != nil {
		var metadata *dto.RequestIncomingMetadata

		if req.Request.Metadata != nil {
			metadata = &dto.RequestIncomingMetadata{
				Body:            req.Request.Metadata.Body,
				Headers:         req.Request.Metadata.Headers,
				QueryParams:     req.Request.Metadata.QueryParams,
				BodyContentType: enums.RequestBodyContentType(req.Request.Metadata.BodyContentType),
			}
		}

		request = &dto.RequestIncoming{
			Path:        req.Request.Path,
			Method:      req.Request.Method,
			Metadata:    metadata,
			IPAddress:   req.Request.IpAddress,
			RequestTime: req.Request.RequestTime.AsTime(),
		}
	}

	return &dto.APIKeyValidate{
		APIKey:         req.ApiKey,
		Request:        request,
		ServiceName:    req.ServiceName,
		ServiceVersion: req.ServiceVersion,
	}
}

func validateResponseFromDomain(response *dto.APIKeyValidateResponse) *pb.ValidateResponse {
	var client *pb.Client
	if response.Client != nil {
		client = &pb.Client{
			Id:   int64(response.Client.ID),
			Name: response.Client.Name,
		}
	}

	var project *pb.Project
	if response.Project != nil {
		project = &pb.Project{
			Id:   int64(response.Project.ID),
			Name: response.Project.Name,
		}
	}

	var environment *pb.Environment
	if response.Environment != nil {
		environment = &pb.Environment{
			Id:   int64(response.Environment.ID),
			Name: response.Environment.Name,
		}
	}

	return &pb.ValidateResponse{
		Valid:       response.Valid,
		RequestId:   response.RequestID,
		FailureCode: string(response.FailureCode),
		Client:      client,
		Project:     project,
		Environment: environment,
	}
}

func validateConsumeResponseFromDomain(
	response *dto.APIKeyValidateConsumeResponse,
) *pb.ValidateConsumeResponse {
	var client *pb.Client
	if response.Client != nil {
		client = &pb.Client{
			Id:   int64(response.Client.ID),
			Name: response.Client.Name,
		}
	}

	var project *pb.Project
	if response.Project != nil {
		project = &pb.Project{
			Id:   int64(response.Project.ID),
			Name: response.Project.Name,
		}
	}

	var environment *pb.Environment
	if response.Environment != nil {
		environment = &pb.Environment{
			Id:   int64(response.Environment.ID),
			Name: response.Environment.Name,
		}
	}

	return &pb.ValidateConsumeResponse{
		BaseResponse: &pb.ValidateResponse{
			Valid:       response.Valid,
			RequestId:   response.RequestID,
			FailureCode: string(response.FailureCode),
			Client:      client,
			Project:     project,
			Environment: environment,
		},
		AvailableRequest: int64(response.AvailableRequest),
	}
}
