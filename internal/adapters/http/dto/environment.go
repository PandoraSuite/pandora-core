package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type EnvironmentService struct {
	ID         int `json:"id" binding:"required"`
	MaxRequest int `json:"max_requests,omitempty"`
}

func (e *EnvironmentService) ToDomain() *dto.EnvironmentService {
	return &dto.EnvironmentService{
		ID:         e.ID,
		MaxRequest: e.MaxRequest,
	}
}

type EnvironmentCreate struct {
	Name      string `json:"name" binding:"required"`
	ProjectID int    `json:"project_id" binding:"required"`

	Services []*EnvironmentService `json:"services" binding:"required"`
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
	Name string `json:"name,omitempty"`
}

func (e *EnvironmentUpdate) ToDomain() *dto.EnvironmentUpdate {
	return &dto.EnvironmentUpdate{
		Name: e.Name,
	}
}

type EnvironmentServiceUpdate struct {
	MaxRequest int `json:"max_requests,omitempty"`
}

func (e *EnvironmentServiceUpdate) ToDomain() *dto.EnvironmentServiceUpdateInput {
	return &dto.EnvironmentServiceUpdateInput{
		MaxRequest: e.MaxRequest,
	}
}

// ... Responses ...

type EnvironmentServiceResponse struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Version          string    `json:"version"`
	MaxRequest       int       `json:"max_requests"`
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
		MaxRequest:       service.MaxRequest,
		AvailableRequest: service.AvailableRequest,
		AssignedAt:       service.AssignedAt,
	}
}

type EnvironmentResponse struct {
	ID        int                           `json:"id"`
	Name      string                        `json:"name"`
	Status    enums.EnvironmentStatus       `json:"status"`
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
		Status:    env.Status,
		ProjectID: env.ProjectID,
		CreatedAt: env.CreatedAt,
		Services:  services,
	}
}

type EnvironmentServiceReset struct {
	ID      int                         `json:"id"`
	Name    string                      `json:"name"`
	Status  enums.EnvironmentStatus     `json:"status"`
	Service *EnvironmentServiceResponse `json:"service"`
}

func EnvironmentServiceResetFromDomain(
	service *dto.EnvironmentServiceReset,
) *EnvironmentServiceReset {
	return &EnvironmentServiceReset{
		ID:      service.ID,
		Name:    service.Name,
		Status:  service.Status,
		Service: EnvironmentServiceResponseFromDomain(service.Service),
	}
}
