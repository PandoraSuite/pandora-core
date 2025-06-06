// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/client/update/ports.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/client/update/ports.go -destination=internal/app/client/update/mock/ports.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	dto "github.com/MAD-py/pandora-core/internal/domain/dto"
	entities "github.com/MAD-py/pandora-core/internal/domain/entities"
	errors "github.com/MAD-py/pandora-core/internal/domain/errors"
	gomock "go.uber.org/mock/gomock"
)

// MockClientRepository is a mock of ClientRepository interface.
type MockClientRepository struct {
	ctrl     *gomock.Controller
	recorder *MockClientRepositoryMockRecorder
	isgomock struct{}
}

// MockClientRepositoryMockRecorder is the mock recorder for MockClientRepository.
type MockClientRepositoryMockRecorder struct {
	mock *MockClientRepository
}

// NewMockClientRepository creates a new mock instance.
func NewMockClientRepository(ctrl *gomock.Controller) *MockClientRepository {
	mock := &MockClientRepository{ctrl: ctrl}
	mock.recorder = &MockClientRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientRepository) EXPECT() *MockClientRepositoryMockRecorder {
	return m.recorder
}

// Update mocks base method.
func (m *MockClientRepository) Update(ctx context.Context, id int, update *dto.ClientUpdate) (*entities.Client, errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, update)
	ret0, _ := ret[0].(*entities.Client)
	ret1, _ := ret[1].(errors.Error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockClientRepositoryMockRecorder) Update(ctx, id, update any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockClientRepository)(nil).Update), ctx, id, update)
}
