package grpc

import (
	"context"
	"strings"

	"github.com/018bf/companies/internal/configs"
	"github.com/018bf/companies/internal/entity"
	"github.com/018bf/companies/pkg/log"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=auth_middleware.go -package=grpc -destination=auth_middleware_mock.go

type ctxKey int

const TokenKey ctxKey = iota + 1

const (
	headerAuthorize = "authorization"
	expectedScheme  = "bearer"
)

func AuthFromMD(ctx context.Context) (string, error) {
	val := metautils.ExtractIncoming(ctx).Get(headerAuthorize)
	if val == "" {
		return "", status.Errorf(
			codes.Unauthenticated,
			"Request unauthenticated with "+expectedScheme,
		)
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", status.Errorf(codes.Unauthenticated, "Bad authorization string")
	}
	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", status.Errorf(
			codes.Unauthenticated,
			"Request unauthenticated with "+expectedScheme,
		)
	}
	return splits[1], nil
}

type authInterceptor interface {
	ValidateToken(ctx context.Context, token *entity.Token) error
}

type AuthMiddleware struct {
	logger          log.Logger
	config          *configs.Config
	authInterceptor authInterceptor
}

func NewAuthMiddleware(
	authInterceptor authInterceptor,
	logger log.Logger,
	config *configs.Config,
) *AuthMiddleware {
	return &AuthMiddleware{authInterceptor: authInterceptor, logger: logger, config: config}
}

func (m *AuthMiddleware) Auth(ctx context.Context) (context.Context, error) {
	var token *entity.Token
	stringToken, err := AuthFromMD(ctx)
	if err != nil {
		return context.WithValue(ctx, TokenKey, token), nil
	}
	if stringToken == "" {
		return context.WithValue(ctx, TokenKey, token), nil
	}
	token = entity.NewToken(stringToken)
	if err := m.authInterceptor.ValidateToken(ctx, token); err != nil {
		return nil, DecodeError(err)
	}
	newCtx := context.WithValue(ctx, TokenKey, token)
	return newCtx, nil
}

func (m *AuthMiddleware) UnaryServerInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	newCtx, err := m.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return handler(newCtx, req)
}
