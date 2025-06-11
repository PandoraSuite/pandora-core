package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type ProjectService struct {
	ID int `json:"id" validate:"required" minimum:"1"`

	MaxRequests int `json:"max_requests" validate:"required" minimum:"-1"`

	ResetFrequency string `json:"reset_frequency" validate:"required" enums:"daily,weekly,biweekly,monthly"`
}

func (p *ProjectService) ToDomain() *dto.ProjectService {
	return &dto.ProjectService{
		ID:             p.ID,
		MaxRequests:    p.MaxRequests,
		ResetFrequency: enums.ProjectServiceResetFrequency(p.ResetFrequency),
	}
}

type ProjectCreate struct {
	Name string `json:"name" validate:"required"`

	ClientID int `json:"client_id" validate:"required" minimum:"1"`

	Services []*ProjectService `json:"services"`
}

func (p *ProjectCreate) ToDomain() *dto.ProjectCreate {
	services := make([]*dto.ProjectService, len(p.Services))
	for i, service := range p.Services {
		services[i] = service.ToDomain()
	}

	return &dto.ProjectCreate{
		Name:     p.Name,
		ClientID: p.ClientID,
		Services: services,
	}
}

type ProjectResetRequest struct {
	RecalculateNextReset bool `json:"recalculate_next_reset" validate:"required"`
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
	NextReset time.Time `json:"next_reset" format:"date"`

	MaxRequests int `json:"max_requests" validate:"required"`

	ResetFrequency string `json:"reset_frequency" enums:"daily,weekly,biweekly,monthly"`
}

func (p *ProjectServiceUpdate) ToDomain() *dto.ProjectServiceUpdate {
	return &dto.ProjectServiceUpdate{
		NextReset:      p.NextReset,
		MaxRequests:    p.MaxRequests,
		ResetFrequency: enums.ProjectServiceResetFrequency(p.ResetFrequency),
	}
}

// ... Reponses ...

type ProjectServiceResponse struct {
	ID int `json:"id" validate:"required" minimum:"1"`

	Name string `json:"name" validate:"required"`

	Version string `json:"version" validate:"required" maxlength:"25"`

	NextReset time.Time `json:"next_reset" validate:"required" format:"date-time" extensions:"x-timezone=utc"`

	MaxRequests int `json:"max_requests" validate:"required" minimum:"-1"`

	ResetFrequency string `json:"reset_frequency" validate:"required" enums:"daily,weekly,biweekly,monthly"`

	AssignedAt time.Time `json:"assigned_at" validate:"required" format:"date-time" extensions:"x-timezone=utc"`
}

func ProjectServiceResponseFromDomain(service *dto.ProjectServiceResponse) *ProjectServiceResponse {
	return &ProjectServiceResponse{
		ID:             service.ID,
		Name:           service.Name,
		Version:        service.Version,
		NextReset:      service.NextReset,
		MaxRequests:    service.MaxRequests,
		ResetFrequency: string(service.ResetFrequency),
		AssignedAt:     service.AssignedAt,
	}
}

type ProjectResponse struct {
	ID int `json:"id" validate:"required" minimum:"1"`

	Name string `json:"name" validate:"required"`

	Status string `json:"status" validate:"required" enums:"enabled,disabled"`

	ClientID int `json:"client_id" validate:"required" minimum:"1"`

	CreatedAt time.Time `json:"created_at" validate:"required" format:"date-time" extensions:"x-timezone=utc"`

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
		Status:    string(project.Status),
		ClientID:  project.ClientID,
		CreatedAt: project.CreatedAt,
		Services:  services,
	}
}

type ProjectResetRequestResponse struct {
	ResetCount int `json:"reset_count" validate:"required" minimum:"0"`

	ProjectService *ProjectServiceResponse `json:"project_service" validate:"required"`

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
