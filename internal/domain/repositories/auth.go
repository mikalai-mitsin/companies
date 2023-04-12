package repositories

import (
	"context"

	"github.com/018bf/companies/internal/domain/models"
)

// AuthRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/auth.go . AuthRepository
type AuthRepository interface {
	Validate(ctx context.Context, token *models.Token) error
	GetSubject(ctx context.Context, token *models.Token) (string, error)
	HasPermission(
		ctx context.Context,
		permission models.PermissionID,
		token *models.Token,
	) error
	HasObjectPermission(
		ctx context.Context,
		permission models.PermissionID,
		token *models.Token,
		obj any,
	) error
}
