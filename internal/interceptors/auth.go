package interceptors

import (
	"context"

	"github.com/018bf/companies/internal/domain/interceptors"
	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/internal/domain/usecases"
	"github.com/018bf/companies/pkg/clock"
	"github.com/018bf/companies/pkg/log"
)

type AuthInterceptor struct {
	authUseCase usecases.AuthUseCase
	clock       clock.Clock
	logger      log.Logger
}

func NewAuthInterceptor(
	authUseCase usecases.AuthUseCase,
	clock clock.Clock,
	logger log.Logger,
) interceptors.AuthInterceptor {
	return &AuthInterceptor{
		authUseCase: authUseCase,
		clock:       clock,
		logger:      logger,
	}
}

func (i *AuthInterceptor) ValidateToken(ctx context.Context, token *models.Token) error {
	if err := i.authUseCase.ValidateToken(ctx, token); err != nil {
		return err
	}
	return nil
}
