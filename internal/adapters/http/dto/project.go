package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type ProjectService struct {
	ID             int                                `json:"id" binding:"required"`
	MaxRequest     int                                `json:"max_requests" binding:"required"`
	ResetFrequency enums.ProjectServiceResetFrequency `json:"reset_frequency"`
}

func (p *ProjectService) ToDomain() *dto.ProjectService {
	return &dto.ProjectService{
		ID:             p.ID,
		MaxRequest:     p.MaxRequest,
		ResetFrequency: p.ResetFrequency,
	}
}

type ProjectCreate struct {
	Name     string              `json:"name" binding:"required"`
	Status   enums.ProjectStatus `json:"status" binding:"required"`
	ClientID int                 `json:"client_id" binding:"required"`

	Services []*ProjectService `json:"services"`
}

func (p *ProjectCreate) ToDomain() *dto.ProjectCreate {
	services := make([]*dto.ProjectService, len(p.Services))
	for i, service := range p.Services {
		services[i] = service.ToDomain()
	}

	return &dto.ProjectCreate{
		Name:     p.Name,
		Status:   p.Status,
		ClientID: p.ClientID,
		Services: services,
	}
}

type ProjectResetRequest struct {
	RecalculateNextReset bool `json:"recalculate_next_reset" binding:"required"`
}

type ProjectUpdate struct {
	Name string `json:"name"`
}

func (p *ProjectUpdate) ToDomain() *dto.ProjectUpdate {
	return &dto.ProjectUpdate{
		Name: p.Name,
	}
}

type ProjectServiceUpdate struct {
	NextReset      time.Time                          `json:"next_reset" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	MaxRequest     int                                `json:"max_requests"`
	ResetFrequency enums.ProjectServiceResetFrequency `json:"reset_frequency"`
}

func (p *ProjectServiceUpdate) ToDomain() *dto.ProjectServiceUpdate {
	return &dto.ProjectServiceUpdate{
		NextReset:      p.NextReset,
		MaxRequest:     p.MaxRequest,
		ResetFrequency: p.ResetFrequency,
	}
}

// ... Reponses ...

type ProjectServiceResponse struct {
	ID             int                                `json:"id"`
	Name           string                             `json:"name"`
	Version        string                             `json:"version"`
	NextReset      time.Time                          `json:"next_reset"`
	MaxRequest     int                                `json:"max_request"`
	ResetFrequency enums.ProjectServiceResetFrequency `json:"reset_frequency"`
	AssignedAt     time.Time                          `json:"assigned_at"`
}

func ProjectServiceResponseFromDomain(service *dto.ProjectServiceResponse) *ProjectServiceResponse {
	return &ProjectServiceResponse{
		ID:             service.ID,
		Name:           service.Name,
		Version:        service.Version,
		NextReset:      service.NextReset,
		MaxRequest:     service.MaxRequest,
		ResetFrequency: service.ResetFrequency,
		AssignedAt:     service.AssignedAt,
	}
}

type ProjectResponse struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	Status    enums.ProjectStatus `json:"status"`
	ClientID  int                 `json:"client_id"`
	CreatedAt time.Time           `json:"created_at"`

	Services []*ProjectServiceResponse `json:"services"`
}

func ProjectResponseFromDomain(project *dto.ProjectResponse) *ProjectResponse {
	services := make([]*ProjectServiceResponse, len(project.Services))
	for i, service := range project.Services {
		services[i] = ProjectServiceResponseFromDomain(service)
	}

	return &ProjectResponse{
		ID:        project.ID,
		Name:      project.Name,
		Status:    project.Status,
		ClientID:  project.ClientID,
		CreatedAt: project.CreatedAt,
		Services:  services,
	}
}

type ProjectResetRequestResponse struct {
	ResetCount          int                        `json:"reset_count"`
	ProjectService      *ProjectServiceResponse    `json:"project_service"`
	EnvironmentServices []*EnvironmentServiceReset `json:"environment_services"`
}

func ProjectResetRequestResponseFromDomain(
	reset *dto.ProjectResetRequestResponse,
) *ProjectResetRequestResponse {
	services := make([]*EnvironmentServiceReset, len(reset.EnvironmentServices))
	for i, service := range reset.EnvironmentServices {
		services[i] = EnvironmentServiceResetFromDomain(service)
	}

	return &ProjectResetRequestResponse{
		ResetCount:          reset.ResetCount,
		ProjectService:      ProjectServiceResponseFromDomain(reset.ProjectService),
		EnvironmentServices: services,
	}
}
