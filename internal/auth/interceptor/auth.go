package interceptor

import (
	"context"

	"github.com/018bf/companies/internal/entity"
	"github.com/018bf/companies/pkg/clock"
	"github.com/018bf/companies/pkg/log"
)

//go:generate mockgen -source=auth.go -package=interceptors -destination=auth_mock.go

type authService interface {
	HasPermission(ctx context.Context, token *entity.Token, permission entity.PermissionID) error
	HasObjectPermission(
		ctx context.Context,
		token *entity.Token,
		permission entity.PermissionID,
		object any,
	) error
	ValidateToken(ctx context.Context, access *entity.Token) error
}

type AuthInterceptor struct {
	authService authService
	clock       clock.Clock
	logger      log.Logger
}

func NewAuthInterceptor(
	authService authService,
	clock clock.Clock,
	logger log.Logger,
) *AuthInterceptor {
	return &AuthInterceptor{
		authService: authService,
		clock:       clock,
		logger:      logger,
	}
}

func (i *AuthInterceptor) ValidateToken(ctx context.Context, token *entity.Token) error {
	if err := i.authService.ValidateToken(ctx, token); err != nil {
		return err
	}
	return nil
}
