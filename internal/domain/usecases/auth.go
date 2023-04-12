package usecases

import (
	"context"

	"github.com/018bf/companies/internal/domain/models"
)

// AuthUseCase - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/auth.go . AuthUseCase
type AuthUseCase interface {
	HasPermission(ctx context.Context, token *models.Token, permission models.PermissionID) error
	HasObjectPermission(
		ctx context.Context,
		token *models.Token,
		permission models.PermissionID,
		object any,
	) error
	ValidateToken(ctx context.Context, access *models.Token) error
}
