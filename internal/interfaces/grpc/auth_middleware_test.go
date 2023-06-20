package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/018bf/companies/pkg/utils"
	"reflect"
	"testing"

	"github.com/018bf/companies/internal/configs"
	"github.com/018bf/companies/internal/entity"
	"github.com/018bf/companies/internal/errs"
	"github.com/018bf/companies/pkg/log"

	"github.com/golang/mock/gomock"

	"google.golang.org/grpc/metadata"
)

type Transport struct {
	MethodPath string
}

func (t Transport) Method() string {
	return t.MethodPath
}

func (t Transport) SetHeader(_ metadata.MD) error {
	return nil
}

func (t Transport) SendHeader(_ metadata.MD) error {
	return nil
}

func (t Transport) SetTrailer(_ metadata.MD) error {
	return nil
}

func TestAuthMiddleware_Auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthInterceptor := NewMockauthInterceptor(ctrl)
	ctx := context.Background()
	token := utils.Pointer(entity.Token("my_token"))
	var tokenPointer *entity.Token
	ctxWithToken := metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", token.String()),
	}))
	ctxWithBadToken := metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", "very bad token"),
	}))
	type fields struct {
		authInterceptor authInterceptor
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    context.Context
		wantErr error
		setup   func()
	}{
		{
			name: "ok",
			setup: func() {
				mockAuthInterceptor.EXPECT().ValidateToken(ctxWithToken, token).Return(nil)
			},
			fields: fields{
				authInterceptor: mockAuthInterceptor,
			},
			args: args{
				ctx: ctxWithToken,
			},
			want:    context.WithValue(ctxWithToken, TokenKey, token),
			wantErr: nil,
		},
		{
			name: "bad token",
			setup: func() {
				mockAuthInterceptor.EXPECT().
					ValidateToken(ctxWithBadToken, entity.NewToken("very bad token")).
					Return(errs.NewBadToken())
			},
			fields: fields{
				authInterceptor: mockAuthInterceptor,
			},
			args: args{
				ctx: ctxWithBadToken,
			},
			want:    nil,
			wantErr: DecodeError(errs.NewBadToken()),
		},
		{
			name:  "without token",
			setup: func() {},
			fields: fields{
				authInterceptor: mockAuthInterceptor,
			},
			args: args{
				ctx: ctx,
			},
			want:    context.WithValue(ctx, TokenKey, tokenPointer),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			m := AuthMiddleware{
				authInterceptor: tt.fields.authInterceptor,
			}
			got, err := m.Auth(tt.args.ctx)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Auth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthInterceptor := NewMockauthInterceptor(ctrl)
	type args struct {
		authInterceptor authInterceptor
		logger          log.Logger
		config          *configs.Config
	}
	tests := []struct {
		name string
		args args
		want *AuthMiddleware
	}{
		{
			name: "ok",
			args: args{
				authInterceptor: mockAuthInterceptor,
			},
			want: &AuthMiddleware{
				authInterceptor: mockAuthInterceptor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthMiddleware(tt.args.authInterceptor, tt.args.logger, tt.args.config); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewAuthMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
