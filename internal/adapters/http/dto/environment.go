package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

// ... Requests ...

type EnvironmentService struct {
	ID          int `json:"id" binding:"required"`
	MaxRequests int `json:"max_requests" binding:"required"`
}

func (e *EnvironmentService) ToDomain() *dto.EnvironmentService {
	return &dto.EnvironmentService{
		ID:          e.ID,
		MaxRequests: e.MaxRequests,
	}
}

type EnvironmentCreate struct {
	Name      string `json:"name" binding:"required"`
	ProjectID int    `json:"project_id" binding:"required"`

	Services []*EnvironmentService `json:"services"`
}

func (e *EnvironmentCreate) ToDomain() *dto.EnvironmentCreate {
	services := make([]*dto.EnvironmentService, len(e.Services))
	for i, service := range e.Services {
		services[i] = service.ToDomain()
	}

	return &dto.EnvironmentCreate{
		Name:      e.Name,
		ProjectID: e.ProjectID,
		Services:  services,
	}
}

type EnvironmentUpdate struct {
	Name string `json:"name"`
}

func (e *EnvironmentUpdate) ToDomain() *dto.EnvironmentUpdate {
	return &dto.EnvironmentUpdate{
		Name: e.Name,
	}
}

type EnvironmentServiceUpdate struct {
	MaxRequests int `json:"max_requests"`
}

func (e *EnvironmentServiceUpdate) ToDomain() *dto.EnvironmentServiceUpdateInput {
	return &dto.EnvironmentServiceUpdateInput{
		MaxRequests: e.MaxRequests,
	}
}

// ... Responses ...

type EnvironmentServiceResponse struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Version          string    `json:"version"`
	MaxRequests      int       `json:"max_requests"`
	AvailableRequest int       `json:"available_requests"`
	AssignedAt       time.Time `json:"assigned_at"`
}

func EnvironmentServiceResponseFromDomain(
	service *dto.EnvironmentServiceResponse,
) *EnvironmentServiceResponse {
	return &EnvironmentServiceResponse{
		ID:               service.ID,
		Name:             service.Name,
		Version:          service.Version,
		MaxRequests:      service.MaxRequests,
		AvailableRequest: service.AvailableRequest,
		AssignedAt:       service.AssignedAt,
	}
}

type EnvironmentResponse struct {
	ID        int                           `json:"id"`
	Name      string                        `json:"name"`
	Status    string                        `json:"status" enums:"enabled,disabled,deprecated"`
	ProjectID int                           `json:"project_id"`
	CreatedAt time.Time                     `json:"created_at"`
	Services  []*EnvironmentServiceResponse `json:"services"`
}

func EnvironmentResponseFromDomain(
	env *dto.EnvironmentResponse,
) *EnvironmentResponse {
	services := make([]*EnvironmentServiceResponse, len(env.Services))
	for i, service := range env.Services {
		services[i] = EnvironmentServiceResponseFromDomain(service)
	}

	return &EnvironmentResponse{
		ID:        env.ID,
		Name:      env.Name,
		Status:    string(env.Status),
		ProjectID: env.ProjectID,
		CreatedAt: env.CreatedAt,
		Services:  services,
	}
}

type EnvironmentServiceReset struct {
	ID      int                         `json:"id"`
	Name    string                      `json:"name"`
	Status  string                      `json:"status" enums:"enabled,disabled,deprecated"`
	Service *EnvironmentServiceResponse `json:"service"`
}

func EnvironmentServiceResetFromDomain(
	service *dto.EnvironmentServiceReset,
) *EnvironmentServiceReset {
	return &EnvironmentServiceReset{
		ID:      service.ID,
		Name:    service.Name,
		Status:  string(service.Status),
		Service: EnvironmentServiceResponseFromDomain(service.Service),
	}
}
