package interceptors

import (
	"context"

	"github.com/018bf/companies/internal/domain/models"
)

// CompanyInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/company.go . CompanyInterceptor
type CompanyInterceptor interface {
	Get(ctx context.Context, id models.UUID, token *models.Token) (*models.Company, error)
	List(
		ctx context.Context,
		filter *models.CompanyFilter,
		token *models.Token,
	) ([]*models.Company, uint64, error)
	Update(
		ctx context.Context,
		update *models.CompanyUpdate,
		token *models.Token,
	) (*models.Company, error)
	Create(
		ctx context.Context,
		create *models.CompanyCreate,
		token *models.Token,
	) (*models.Company, error)
	Delete(ctx context.Context, id models.UUID, token *models.Token) error
}
