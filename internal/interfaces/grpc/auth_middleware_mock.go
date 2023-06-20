// Code generated by MockGen. DO NOT EDIT.
// Source: auth_middleware.go

// Package grpc is a generated GoMock package.
package grpc

import (
	context "context"
	reflect "reflect"

	models "github.com/018bf/companies/internal/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockauthInterceptor is a mock of authInterceptor interface.
type MockauthInterceptor struct {
	ctrl     *gomock.Controller
	recorder *MockauthInterceptorMockRecorder
}

// MockauthInterceptorMockRecorder is the mock recorder for MockauthInterceptor.
type MockauthInterceptorMockRecorder struct {
	mock *MockauthInterceptor
}

// NewMockauthInterceptor creates a new mock instance.
func NewMockauthInterceptor(ctrl *gomock.Controller) *MockauthInterceptor {
	mock := &MockauthInterceptor{ctrl: ctrl}
	mock.recorder = &MockauthInterceptorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockauthInterceptor) EXPECT() *MockauthInterceptorMockRecorder {
	return m.recorder
}

// ValidateToken mocks base method.
func (m *MockauthInterceptor) ValidateToken(ctx context.Context, token *models.Token) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockauthInterceptorMockRecorder) ValidateToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockauthInterceptor)(nil).ValidateToken), ctx, token)
}
