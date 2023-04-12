package usecases

import (
	"context"

	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/internal/domain/repositories"
	"github.com/018bf/companies/internal/domain/usecases"
	"github.com/018bf/companies/pkg/log"
)

type AuthUseCase struct {
	authRepository repositories.AuthRepository
	logger         log.Logger
}

func NewAuthUseCase(
	authRepository repositories.AuthRepository,
	logger log.Logger,
) usecases.AuthUseCase {
	return &AuthUseCase{
		authRepository: authRepository,
		logger:         logger,
	}
}

func (u AuthUseCase) ValidateToken(ctx context.Context, token *models.Token) error {
	if err := u.authRepository.Validate(ctx, token); err != nil {
		return err
	}
	return nil
}

func (u AuthUseCase) HasPermission(
	ctx context.Context,
	token *models.Token,
	permission models.PermissionID,
) error {
	if err := u.authRepository.HasPermission(ctx, permission, token); err != nil {
		return err
	}
	return nil
}

func (u AuthUseCase) HasObjectPermission(
	ctx context.Context,
	token *models.Token,
	permission models.PermissionID,
	object any,
) error {
	if err := u.authRepository.HasObjectPermission(ctx, permission, token, object); err != nil {
		return err
	}
	return nil
}
