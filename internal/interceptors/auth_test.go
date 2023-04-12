package interceptors

import (
	"context"
	"errors"
	"github.com/018bf/companies/pkg/utils"
	"reflect"
	"testing"

	"github.com/018bf/companies/internal/domain/errs"
	"github.com/018bf/companies/internal/domain/interceptors"
	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/internal/domain/usecases"
	mock_usecases "github.com/018bf/companies/internal/domain/usecases/mock"
	"github.com/018bf/companies/pkg/clock"
	mock_clock "github.com/018bf/companies/pkg/clock/mock"
	"github.com/018bf/companies/pkg/log"
	mock_log "github.com/018bf/companies/pkg/log/mock"
	"github.com/golang/mock/gomock"
)

func TestNewAuthInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	clockmock := mock_clock.NewMockClock(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase usecases.AuthUseCase
		logger      log.Logger
		clock       clock.Clock
	}
	tests := []struct {
		name string
		args args
		want interceptors.AuthInterceptor
	}{
		{
			name: "ok",
			args: args{
				authUseCase: authUseCase,
				logger:      logger,
				clock:       clockmock,
			},
			want: &AuthInterceptor{
				authUseCase: authUseCase,
				clock:       clockmock,
				logger:      logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAuthInterceptor(
				tt.args.authUseCase,
				tt.args.clock,
				tt.args.logger,
			)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthInterceptor_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	token := utils.Pointer(models.Token("this_is_valid_token"))
	type fields struct {
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx   context.Context
		token *models.Token
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
				authUseCase.EXPECT().ValidateToken(ctx, token).Return(nil).Times(1)
			},
			fields: fields{
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:   ctx,
				token: token,
			},
			wantErr: nil,
		},
		{
			name: "repository error",
			setup: func() {
				authUseCase.EXPECT().
					ValidateToken(ctx, token).
					Return(errs.NewUnexpectedBehaviorError("35124345")).
					Times(1)
			},
			fields: fields{
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:   ctx,
				token: token,
			},
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"details": "35124345",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := AuthInterceptor{
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			if err := i.ValidateToken(tt.args.ctx, tt.args.token); !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
