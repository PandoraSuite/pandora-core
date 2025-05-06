package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

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
