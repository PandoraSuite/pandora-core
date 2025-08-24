package resetduerequests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/MAD-py/pandora-core/internal/app/project/reset_due_requests/mock"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/utils"
)

type UseCaseSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	projectRepo *mock.MockProjectRepository

	useCase UseCase

	ctx context.Context
}

func (s *UseCaseSuite) SetupTest() {
	time.Local = time.UTC
	s.ctrl = gomock.NewController(s.T())

	s.projectRepo = mock.NewMockProjectRepository(s.ctrl)

	s.useCase = NewUseCase(s.projectRepo)

	s.ctx = context.Background()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()
}

// Success scenarios

func (s *UseCaseSuite) TestExecute_Success_SingleProjectSingleService() {
	// Mock time for consistent testing
	testTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Test data
	project := &entities.Project{
		ID:       1,
		Name:     "Test Project",
		Status:   enums.ProjectStatusEnabled,
		ClientID: 10,
		Services: []*entities.ProjectService{
			{
				ID:             100,
				Name:           "TestService",
				Version:        "1.0.0",
				NextReset:      testTime,
				MaxRequests:    1000,
				ResetFrequency: enums.ProjectServiceResetFrequencyDaily,
				AssignedAt:     testTime.AddDate(0, 0, -1),
			},
		},
		CreatedAt: testTime.AddDate(0, 0, -30),
	}

	envServiceResets := []*dto.EnvironmentServiceReset{
		{
			ID:     1,
			Name:   "production",
			Status: enums.EnvironmentStatusEnabled,
			Service: &dto.EnvironmentServiceResponse{
				ID:               100,
				Name:             "TestService",
				Version:          "1.0.0",
				MaxRequests:      1000,
				AvailableRequest: 1000,
				AssignedAt:       testTime.AddDate(0, 0, -1),
			},
		},
	}

	// Mock expectations
	s.projectRepo.EXPECT().
		ListProjectServiceDueForReset(s.ctx, utils.TruncateToDay(time.Now())).
		Return([]*entities.Project{project}, nil).
		Times(1)

	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(s.ctx, project.ID, project.Services[0].ID, gomock.Any()).
		DoAndReturn(func(ctx context.Context, projectID, serviceID int, nextResetTime time.Time) ([]*dto.EnvironmentServiceReset, errors.Error) {
			// Verify that CalculateNextServicesReset was called by checking the nextReset time is updated
			s.Require().True(nextResetTime.After(testTime))
			return envServiceResets, nil
		}).
		Times(1)

	// Execute
	result, err := s.useCase.Execute(s.ctx)

	// Assertions
	s.Require().NoError(err)
	s.Require().Len(result, 1)
	
	projectReset := result[0]
	s.Equal(project.ID, projectReset.ID)
	s.Equal(project.Name, projectReset.Name)
	s.Equal(project.Status, projectReset.Status)
	s.Require().Len(projectReset.EnvironmentServices, 1)
	s.Equal(envServiceResets[0], projectReset.EnvironmentServices[0])
}

func (s *UseCaseSuite) TestExecute_Success_MultipleProjectsMultipleServices() {
	testTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Test data with two projects, each having multiple services
	projects := []*entities.Project{
		{
			ID:       1,
			Name:     "Project One",
			Status:   enums.ProjectStatusEnabled,
			ClientID: 10,
			Services: []*entities.ProjectService{
				{
					ID:             100,
					Name:           "ServiceA",
					Version:        "1.0.0",
					ResetFrequency: enums.ProjectServiceResetFrequencyDaily,
				},
				{
					ID:             101,
					Name:           "ServiceB",
					Version:        "2.0.0",
					ResetFrequency: enums.ProjectServiceResetFrequencyWeekly,
				},
			},
			CreatedAt: testTime.AddDate(0, 0, -30),
		},
		{
			ID:       2,
			Name:     "Project Two",
			Status:   enums.ProjectStatusEnabled,
			ClientID: 20,
			Services: []*entities.ProjectService{
				{
					ID:             200,
					Name:           "ServiceC",
					Version:        "1.0.0",
					ResetFrequency: enums.ProjectServiceResetFrequencyMonthly,
				},
			},
			CreatedAt: testTime.AddDate(0, 0, -60),
		},
	}

	envServiceResets1 := []*dto.EnvironmentServiceReset{
		{ID: 1, Name: "env1", Status: enums.EnvironmentStatusEnabled},
	}
	envServiceResets2 := []*dto.EnvironmentServiceReset{
		{ID: 2, Name: "env2", Status: enums.EnvironmentStatusEnabled},
	}
	envServiceResets3 := []*dto.EnvironmentServiceReset{
		{ID: 3, Name: "env3", Status: enums.EnvironmentStatusEnabled},
	}

	// Mock expectations
	s.projectRepo.EXPECT().
		ListProjectServiceDueForReset(s.ctx, utils.TruncateToDay(time.Now())).
		Return(projects, nil).
		Times(1)

	// Project 1 service resets
	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(s.ctx, 1, 100, gomock.Any()).
		Return(envServiceResets1, nil).
		Times(1)

	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(s.ctx, 1, 101, gomock.Any()).
		Return(envServiceResets2, nil).
		Times(1)

	// Project 2 service reset
	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(s.ctx, 2, 200, gomock.Any()).
		Return(envServiceResets3, nil).
		Times(1)

	// Execute
	result, err := s.useCase.Execute(s.ctx)

	// Assertions
	s.Require().NoError(err)
	s.Require().Len(result, 2)

	// Verify first project
	project1Result := result[0]
	s.Equal(1, project1Result.ID)
	s.Equal("Project One", project1Result.Name)
	s.Equal(enums.ProjectStatusEnabled, project1Result.Status)
	s.Require().Len(project1Result.EnvironmentServices, 2)

	// Verify second project
	project2Result := result[1]
	s.Equal(2, project2Result.ID)
	s.Equal("Project Two", project2Result.Name)
	s.Equal(enums.ProjectStatusEnabled, project2Result.Status)
	s.Require().Len(project2Result.EnvironmentServices, 1)
}

func (s *UseCaseSuite) TestExecute_Success_EmptyResults() {
	// Mock expectations - no projects due for reset
	s.projectRepo.EXPECT().
		ListProjectServiceDueForReset(s.ctx, utils.TruncateToDay(time.Now())).
		Return([]*entities.Project{}, nil).
		Times(1)

	// ResetProjectServiceUsage should NOT be called when no projects are found
	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	// Execute
	result, err := s.useCase.Execute(s.ctx)

	// Assertions
	s.Require().NoError(err)
	s.Require().Empty(result)
}

// Error scenarios

func (s *UseCaseSuite) TestExecute_Error_ListProjectServiceDueForResetFails() {
	expectedError := errors.NewInternal("database connection failed", nil)

	// Mock expectations
	s.projectRepo.EXPECT().
		ListProjectServiceDueForReset(s.ctx, utils.TruncateToDay(time.Now())).
		Return(nil, expectedError).
		Times(1)

	// ResetProjectServiceUsage should NOT be called when ListProjectServiceDueForReset fails
	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	// Execute
	result, err := s.useCase.Execute(s.ctx)

	// Assertions
	s.Require().Error(err)
	s.Require().Nil(result)
	s.Equal(expectedError, err)
}

func (s *UseCaseSuite) TestExecute_Error_ResetProjectServiceUsageFails_SingleProject() {
	testTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	resetError := errors.NewInternal("failed to reset service usage", nil)

	// Test data with single project and single service
	project := &entities.Project{
		ID:       1,
		Name:     "Test Project",
		Status:   enums.ProjectStatusEnabled,
		ClientID: 10,
		Services: []*entities.ProjectService{
			{
				ID:             100,
				Name:           "TestService",
				Version:        "1.0.0",
				ResetFrequency: enums.ProjectServiceResetFrequencyDaily,
			},
		},
		CreatedAt: testTime.AddDate(0, 0, -30),
	}

	// Mock expectations
	s.projectRepo.EXPECT().
		ListProjectServiceDueForReset(s.ctx, utils.TruncateToDay(time.Now())).
		Return([]*entities.Project{project}, nil).
		Times(1)

	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(s.ctx, project.ID, project.Services[0].ID, gomock.Any()).
		Return(nil, resetError).
		Times(1)

	// Execute
	result, err := s.useCase.Execute(s.ctx)

	// Assertions
	// The use case returns partial results even when some operations fail
	// but does not return an aggregate error in this implementation
	s.Require().NoError(err)
	s.Require().Len(result, 1)
	
	projectResult := result[0]
	s.Equal(project.ID, projectResult.ID)
	s.Equal(project.Name, projectResult.Name)
	s.Equal(project.Status, projectResult.Status)
	// Environment services should be empty since the reset failed
	s.Empty(projectResult.EnvironmentServices)
}

func (s *UseCaseSuite) TestExecute_Error_PartialFailure_MultipleProjects() {
	testTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	resetError := errors.NewInternal("failed to reset service", nil)

	// Test data with two projects
	projects := []*entities.Project{
		{
			ID:       1,
			Name:     "Project One",
			Status:   enums.ProjectStatusEnabled,
			ClientID: 10,
			Services: []*entities.ProjectService{
				{
					ID:             100,
					Name:           "ServiceA",
					Version:        "1.0.0",
					ResetFrequency: enums.ProjectServiceResetFrequencyDaily,
				},
				{
					ID:             101,
					Name:           "ServiceB",
					Version:        "2.0.0",
					ResetFrequency: enums.ProjectServiceResetFrequencyWeekly,
				},
			},
			CreatedAt: testTime.AddDate(0, 0, -30),
		},
		{
			ID:       2,
			Name:     "Project Two",
			Status:   enums.ProjectStatusEnabled,
			ClientID: 20,
			Services: []*entities.ProjectService{
				{
					ID:             200,
					Name:           "ServiceC",
					Version:        "1.0.0",
					ResetFrequency: enums.ProjectServiceResetFrequencyMonthly,
				},
			},
			CreatedAt: testTime.AddDate(0, 0, -60),
		},
	}

	envServiceResets := []*dto.EnvironmentServiceReset{
		{ID: 1, Name: "env1", Status: enums.EnvironmentStatusEnabled},
	}

	// Mock expectations
	s.projectRepo.EXPECT().
		ListProjectServiceDueForReset(s.ctx, utils.TruncateToDay(time.Now())).
		Return(projects, nil).
		Times(1)

	// Project 1 - first service succeeds, second fails
	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(s.ctx, 1, 100, gomock.Any()).
		Return(envServiceResets, nil).
		Times(1)

	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(s.ctx, 1, 101, gomock.Any()).
		Return(nil, resetError).
		Times(1)

	// Project 2 - service succeeds
	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(s.ctx, 2, 200, gomock.Any()).
		Return(envServiceResets, nil).
		Times(1)

	// Execute
	result, err := s.useCase.Execute(s.ctx)

	// Assertions
	// The current implementation continues processing and doesn't return aggregate errors
	s.Require().NoError(err)
	s.Require().Len(result, 2)

	// Verify first project - should have results from successful service only
	project1Result := result[0]
	s.Equal(1, project1Result.ID)
	s.Equal("Project One", project1Result.Name)
	s.Require().Len(project1Result.EnvironmentServices, 1) // Only successful service

	// Verify second project - should have results from successful service
	project2Result := result[1]
	s.Equal(2, project2Result.ID)
	s.Equal("Project Two", project2Result.Name)
	s.Require().Len(project2Result.EnvironmentServices, 1)
}

func (s *UseCaseSuite) TestExecute_Error_AllServiceResetsFail() {
	testTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	resetError := errors.NewInternal("failed to reset all services", nil)

	// Test data with single project and single service
	project := &entities.Project{
		ID:       1,
		Name:     "Test Project",
		Status:   enums.ProjectStatusEnabled,
		ClientID: 10,
		Services: []*entities.ProjectService{
			{
				ID:             100,
				Name:           "TestService",
					Version:        "1.0.0",
				ResetFrequency: enums.ProjectServiceResetFrequencyDaily,
			},
		},
		CreatedAt: testTime.AddDate(0, 0, -30),
	}

	// Mock expectations
	s.projectRepo.EXPECT().
		ListProjectServiceDueForReset(s.ctx, utils.TruncateToDay(time.Now())).
		Return([]*entities.Project{project}, nil).
		Times(1)

	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(s.ctx, project.ID, project.Services[0].ID, gomock.Any()).
		Return(nil, resetError).
		Times(1)

	// Execute
	result, err := s.useCase.Execute(s.ctx)

	// Assertions
	s.Require().NoError(err) // Current implementation doesn't return errors from individual resets
	s.Require().Len(result, 1)
	
	projectResult := result[0]
	s.Equal(project.ID, projectResult.ID)
	s.Equal(project.Name, projectResult.Name)
	s.Equal(project.Status, projectResult.Status)
	// Environment services should be empty since all resets failed
	s.Empty(projectResult.EnvironmentServices)
}

// Edge cases

func (s *UseCaseSuite) TestExecute_EdgeCase_ProjectWithNoServices() {
	testTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Test data with project having no services
	project := &entities.Project{
		ID:        1,
		Name:      "Empty Project",
		Status:    enums.ProjectStatusEnabled,
		ClientID:  10,
		Services:  []*entities.ProjectService{}, // Empty services
		CreatedAt: testTime.AddDate(0, 0, -30),
	}

	// Mock expectations
	s.projectRepo.EXPECT().
		ListProjectServiceDueForReset(s.ctx, utils.TruncateToDay(time.Now())).
		Return([]*entities.Project{project}, nil).
		Times(1)

	// ResetProjectServiceUsage should NOT be called when project has no services
	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	// Execute
	result, err := s.useCase.Execute(s.ctx)

	// Assertions
	s.Require().NoError(err)
	s.Require().Len(result, 1)
	
	projectResult := result[0]
	s.Equal(project.ID, projectResult.ID)
	s.Equal(project.Name, projectResult.Name)
	s.Equal(project.Status, projectResult.Status)
	s.Empty(projectResult.EnvironmentServices)
}

func (s *UseCaseSuite) TestExecute_EdgeCase_ContextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context before execution

	// Mock expectations - the repository call should still be attempted
	s.projectRepo.EXPECT().
		ListProjectServiceDueForReset(ctx, utils.TruncateToDay(time.Now())).
		Return(nil, errors.NewInternal("context canceled", context.Canceled)).
		Times(1)

	// ResetProjectServiceUsage should NOT be called when ListProjectServiceDueForReset fails
	s.projectRepo.EXPECT().
		ResetProjectServiceUsage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	// Execute with canceled context
	result, err := s.useCase.Execute(ctx)

	// Assertions
	s.Require().Error(err)
	s.Require().Nil(result)
	s.Equal(errors.CodeInternal, err.Code())
}

func TestUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UseCaseSuite))
}
