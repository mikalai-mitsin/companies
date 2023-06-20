// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go

// Package usecases is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	models "github.com/018bf/companies/internal/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockauthRepository is a mock of authRepository interface.
type MockauthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockauthRepositoryMockRecorder
}

// MockauthRepositoryMockRecorder is the mock recorder for MockauthRepository.
type MockauthRepositoryMockRecorder struct {
	mock *MockauthRepository
}

// NewMockauthRepository creates a new mock instance.
func NewMockauthRepository(ctrl *gomock.Controller) *MockauthRepository {
	mock := &MockauthRepository{ctrl: ctrl}
	mock.recorder = &MockauthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockauthRepository) EXPECT() *MockauthRepositoryMockRecorder {
	return m.recorder
}

// GetSubject mocks base method.
func (m *MockauthRepository) GetSubject(ctx context.Context, token *models.Token) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubject", ctx, token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubject indicates an expected call of GetSubject.
func (mr *MockauthRepositoryMockRecorder) GetSubject(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubject", reflect.TypeOf((*MockauthRepository)(nil).GetSubject), ctx, token)
}

// HasObjectPermission mocks base method.
func (m *MockauthRepository) HasObjectPermission(ctx context.Context, permission models.PermissionID, token *models.Token, obj any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasObjectPermission", ctx, permission, token, obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// HasObjectPermission indicates an expected call of HasObjectPermission.
func (mr *MockauthRepositoryMockRecorder) HasObjectPermission(ctx, permission, token, obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasObjectPermission", reflect.TypeOf((*MockauthRepository)(nil).HasObjectPermission), ctx, permission, token, obj)
}

// HasPermission mocks base method.
func (m *MockauthRepository) HasPermission(ctx context.Context, permission models.PermissionID, token *models.Token) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasPermission", ctx, permission, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// HasPermission indicates an expected call of HasPermission.
func (mr *MockauthRepositoryMockRecorder) HasPermission(ctx, permission, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasPermission", reflect.TypeOf((*MockauthRepository)(nil).HasPermission), ctx, permission, token)
}

// Validate mocks base method.
func (m *MockauthRepository) Validate(ctx context.Context, token *models.Token) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockauthRepositoryMockRecorder) Validate(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockauthRepository)(nil).Validate), ctx, token)
}
