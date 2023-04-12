package interceptors

import (
	"context"

	"github.com/018bf/companies/internal/domain/models"
)

// AuthInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/auth.go . AuthInterceptor
type AuthInterceptor interface {
	ValidateToken(ctx context.Context, token *models.Token) error
}
