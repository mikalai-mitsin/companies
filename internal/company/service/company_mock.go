// Code generated by MockGen. DO NOT EDIT.
// Source: company.go

// Package usecases is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	models "github.com/018bf/companies/internal/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockcompanyRepository is a mock of companyRepository interface.
type MockcompanyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockcompanyRepositoryMockRecorder
}

// MockcompanyRepositoryMockRecorder is the mock recorder for MockcompanyRepository.
type MockcompanyRepositoryMockRecorder struct {
	mock *MockcompanyRepository
}

// NewMockcompanyRepository creates a new mock instance.
func NewMockcompanyRepository(ctrl *gomock.Controller) *MockcompanyRepository {
	mock := &MockcompanyRepository{ctrl: ctrl}
	mock.recorder = &MockcompanyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcompanyRepository) EXPECT() *MockcompanyRepositoryMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockcompanyRepository) Count(ctx context.Context, filter *models.CompanyFilter) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, filter)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockcompanyRepositoryMockRecorder) Count(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockcompanyRepository)(nil).Count), ctx, filter)
}

// Create mocks base method.
func (m *MockcompanyRepository) Create(ctx context.Context, create *models.Company) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, create)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockcompanyRepositoryMockRecorder) Create(ctx, create interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockcompanyRepository)(nil).Create), ctx, create)
}

// Delete mocks base method.
func (m *MockcompanyRepository) Delete(ctx context.Context, id models.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockcompanyRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockcompanyRepository)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockcompanyRepository) Get(ctx context.Context, id models.UUID) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockcompanyRepositoryMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockcompanyRepository)(nil).Get), ctx, id)
}

// List mocks base method.
func (m *MockcompanyRepository) List(ctx context.Context, filter *models.CompanyFilter) ([]*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filter)
	ret0, _ := ret[0].([]*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockcompanyRepositoryMockRecorder) List(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockcompanyRepository)(nil).List), ctx, filter)
}

// Update mocks base method.
func (m *MockcompanyRepository) Update(ctx context.Context, update *models.Company) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, update)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockcompanyRepositoryMockRecorder) Update(ctx, update interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockcompanyRepository)(nil).Update), ctx, update)
}