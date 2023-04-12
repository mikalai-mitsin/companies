package usecases

import (
	"context"
	"errors"
	"github.com/018bf/companies/pkg/utils"
	"reflect"
	"testing"

	"github.com/018bf/companies/internal/domain/errs"
	"github.com/018bf/companies/internal/domain/models"
	mock_models "github.com/018bf/companies/internal/domain/models/mock"
	"github.com/018bf/companies/internal/domain/repositories"
	mock_repositories "github.com/018bf/companies/internal/domain/repositories/mock"
	"github.com/018bf/companies/internal/domain/usecases"
	"github.com/018bf/companies/pkg/log"
	mock_log "github.com/018bf/companies/pkg/log/mock"

	"github.com/golang/mock/gomock"
)

func TestAuthUseCase_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authRepository := mock_repositories.NewMockAuthRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	type fields struct {
		authRepository repositories.AuthRepository
		logger         log.Logger
	}
	type args struct {
		ctx   context.Context
		token *models.Token
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
				authRepository.EXPECT().Validate(ctx, utils.Pointer(models.Token("my_token"))).Return(nil).Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				logger:         logger,
			},
			args: args{
				ctx:   ctx,
				token: utils.Pointer(models.Token("my_token")),
			},
			wantErr: nil,
		},
		{
			name: "repository error",
			setup: func() {
				authRepository.EXPECT().
					Validate(ctx, utils.Pointer(models.Token("my_token"))).
					Return(errs.NewUnexpectedBehaviorError("error 345")).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				logger:         logger,
			},
			args: args{
				ctx:   ctx,
				token: utils.Pointer(models.Token("my_token")),
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
			u := AuthUseCase{
				authRepository: tt.fields.authRepository,
				logger:         tt.fields.logger,
			}
			if err := u.ValidateToken(tt.args.ctx, tt.args.token); !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAuthUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authRepository := mock_repositories.NewMockAuthRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authRepository repositories.AuthRepository
		logger         log.Logger
	}
	tests := []struct {
		name string
		args args
		want usecases.AuthUseCase
	}{
		{
			name: "ok",
			args: args{
				authRepository: authRepository,
				logger:         logger,
			},
			want: &AuthUseCase{
				authRepository: authRepository,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthUseCase(tt.args.authRepository, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewAuthUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthUseCase_HasPermission(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authRepository := mock_repositories.NewMockAuthRepository(ctrl)
	user := utils.Pointer(mock_models.NewToken(t))
	type fields struct {
		authRepository repositories.AuthRepository
		logger         log.Logger
	}
	type args struct {
		in0        context.Context
		in1        *models.Token
		permission models.PermissionID
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
				authRepository.EXPECT().
					HasPermission(ctx, models.PermissionIDCompanyDelete, user).
					Return(nil)
			},
			fields: fields{
				authRepository: authRepository,
				logger:         nil,
			},
			args: args{
				in0:        ctx,
				in1:        user,
				permission: models.PermissionIDCompanyDelete,
			},
			wantErr: nil,
		},
		{
			name: "error",
			setup: func() {
				authRepository.EXPECT().
					HasPermission(ctx, models.PermissionIDCompanyDelete, user).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authRepository: authRepository,
				logger:         nil,
			},
			args: args{
				in0:        ctx,
				in1:        user,
				permission: models.PermissionIDCompanyDelete,
			},
			wantErr: errs.NewPermissionDenied(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := AuthUseCase{
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

func TestAuthUseCase_HasObjectPermission(t *testing.T) {
	ctx := context.Background()
	user := utils.Pointer(mock_models.NewToken(t))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authRepository := mock_repositories.NewMockAuthRepository(ctrl)
	type fields struct {
		authRepository repositories.AuthRepository
		logger         log.Logger
	}
	type args struct {
		in0        context.Context
		user       *models.Token
		permission models.PermissionID
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
				authRepository.EXPECT().
					HasObjectPermission(ctx, models.PermissionIDCompanyDetail, user, "user").
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authRepository: authRepository,
				logger:         nil,
			},
			args: args{
				in0:        ctx,
				user:       user,
				permission: models.PermissionIDCompanyDetail,
				object:     "user",
			},
			wantErr: errs.NewPermissionDenied(),
		},
		{
			name: "ok",
			setup: func() {
				authRepository.EXPECT().
					HasObjectPermission(ctx, models.PermissionIDCompanyDetail, user, user).
					Return(nil)
			},
			fields: fields{
				authRepository: authRepository,
				logger:         nil,
			},
			args: args{
				in0:        ctx,
				user:       user,
				permission: models.PermissionIDCompanyDetail,
				object:     user,
			},
			wantErr: nil,
		},
		{
			name: "ok with user",
			setup: func() {
				authRepository.EXPECT().
					HasObjectPermission(ctx, models.PermissionIDCompanyList, user, user).
					Return(nil)
			},
			fields: fields{
				authRepository: authRepository,
				logger:         nil,
			},
			args: args{
				in0:        ctx,
				user:       user,
				permission: models.PermissionIDCompanyList,
				object:     user,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := AuthUseCase{
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
