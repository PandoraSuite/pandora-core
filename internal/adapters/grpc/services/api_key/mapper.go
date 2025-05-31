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
	}
}

func validateConsumeResponseFromDomain(
	response *dto.APIKeyValidateConsumeResponse,
) *pb.ValidateConsumeResponse {
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
	}
}
