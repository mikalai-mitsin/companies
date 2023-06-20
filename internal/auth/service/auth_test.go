package service

import (
	"context"
	"errors"
	"github.com/018bf/companies/pkg/utils"
	"reflect"
	"testing"

	"github.com/018bf/companies/internal/entity"
	mock_models "github.com/018bf/companies/internal/entity/mock"
	"github.com/018bf/companies/internal/errs"
	"github.com/018bf/companies/pkg/log"
	mock_log "github.com/018bf/companies/pkg/log/mock"

	"github.com/golang/mock/gomock"
)

func TestAuthService_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepository := NewMockauthRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	type fields struct {
		authRepository authRepository
		logger         log.Logger
	}
	type args struct {
		ctx   context.Context
		token *entity.Token
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		setup   func()
	}{
		{
			name: "ok",
			setup: func() {
				mockAuthRepository.EXPECT().Validate(ctx, utils.Pointer(entity.Token("my_token"))).Return(nil).Times(1)
			},
			fields: fields{
				authRepository: mockAuthRepository,
				logger:         logger,
			},
			args: args{
				ctx:   ctx,
				token: utils.Pointer(entity.Token("my_token")),
			},
			wantErr: nil,
		},
		{
			name: "repository error",
			setup: func() {
				mockAuthRepository.EXPECT().
					Validate(ctx, utils.Pointer(entity.Token("my_token"))).
					Return(errs.NewUnexpectedBehaviorError("error 345")).
					Times(1)
			},
			fields: fields{
				authRepository: mockAuthRepository,
				logger:         logger,
			},
			args: args{
				ctx:   ctx,
				token: utils.Pointer(entity.Token("my_token")),
			},
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params:  map[string]string{"details": "error 345"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := AuthService{
				authRepository: tt.fields.authRepository,
				logger:         tt.fields.logger,
			}
			if err := u.ValidateToken(tt.args.ctx, tt.args.token); !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAuthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepository := NewMockauthRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authRepository authRepository
		logger         log.Logger
	}
	tests := []struct {
		name string
		args args
		want *AuthService
	}{
		{
			name: "ok",
			args: args{
				authRepository: mockAuthRepository,
				logger:         logger,
			},
			want: &AuthService{
				authRepository: mockAuthRepository,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthService(tt.args.authRepository, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewAuthService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_HasPermission(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepository := NewMockauthRepository(ctrl)
	user := utils.Pointer(mock_models.NewToken(t))
	type fields struct {
		authRepository authRepository
		logger         log.Logger
	}
	type args struct {
		in0        context.Context
		in1        *entity.Token
		permission entity.PermissionID
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockAuthRepository.EXPECT().
					HasPermission(ctx, entity.PermissionIDCompanyDelete, user).
					Return(nil)
			},
			fields: fields{
				authRepository: mockAuthRepository,
				logger:         nil,
			},
			args: args{
				in0:        ctx,
				in1:        user,
				permission: entity.PermissionIDCompanyDelete,
			},
			wantErr: nil,
		},
		{
			name: "error",
			setup: func() {
				mockAuthRepository.EXPECT().
					HasPermission(ctx, entity.PermissionIDCompanyDelete, user).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authRepository: mockAuthRepository,
				logger:         nil,
			},
			args: args{
				in0:        ctx,
				in1:        user,
				permission: entity.PermissionIDCompanyDelete,
			},
			wantErr: errs.NewPermissionDenied(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := AuthService{
				authRepository: tt.fields.authRepository,
				logger:         tt.fields.logger,
			}
			tt.setup()
			if err := u.HasPermission(tt.args.in0, tt.args.in1, tt.args.permission); !errors.Is(
				err,
				tt.wantErr,
			) {
				t.Errorf("HasPermission() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_HasObjectPermission(t *testing.T) {
	ctx := context.Background()
	user := utils.Pointer(mock_models.NewToken(t))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepository := NewMockauthRepository(ctrl)
	type fields struct {
		authRepository authRepository
		logger         log.Logger
	}
	type args struct {
		in0        context.Context
		user       *entity.Token
		permission entity.PermissionID
		object     interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		setup   func()
	}{
		{
			name: "error",
			setup: func() {
				mockAuthRepository.EXPECT().
					HasObjectPermission(ctx, entity.PermissionIDCompanyDetail, user, "user").
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authRepository: mockAuthRepository,
				logger:         nil,
			},
			args: args{
				in0:        ctx,
				user:       user,
				permission: entity.PermissionIDCompanyDetail,
				object:     "user",
			},
			wantErr: errs.NewPermissionDenied(),
		},
		{
			name: "ok",
			setup: func() {
				mockAuthRepository.EXPECT().
					HasObjectPermission(ctx, entity.PermissionIDCompanyDetail, user, user).
					Return(nil)
			},
			fields: fields{
				authRepository: mockAuthRepository,
				logger:         nil,
			},
			args: args{
				in0:        ctx,
				user:       user,
				permission: entity.PermissionIDCompanyDetail,
				object:     user,
			},
			wantErr: nil,
		},
		{
			name: "ok with user",
			setup: func() {
				mockAuthRepository.EXPECT().
					HasObjectPermission(ctx, entity.PermissionIDCompanyList, user, user).
					Return(nil)
			},
			fields: fields{
				authRepository: mockAuthRepository,
				logger:         nil,
			},
			args: args{
				in0:        ctx,
				user:       user,
				permission: entity.PermissionIDCompanyList,
				object:     user,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := AuthService{
				authRepository: tt.fields.authRepository,
				logger:         tt.fields.logger,
			}
			tt.setup()
			if err := u.HasObjectPermission(tt.args.in0, tt.args.user, tt.args.permission, tt.args.object); !errors.Is(
				err,
				tt.wantErr,
			) {
				t.Errorf("HasObjectPermission() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
