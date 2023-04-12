package repositories

import (
	"context"

	"github.com/018bf/companies/internal/domain/models"
)

// CompanyRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/company.go . CompanyRepository
type CompanyRepository interface {
	Get(ctx context.Context, id models.UUID) (*models.Company, error)
	List(ctx context.Context, filter *models.CompanyFilter) ([]*models.Company, error)
	Count(ctx context.Context, filter *models.CompanyFilter) (uint64, error)
	Update(ctx context.Context, update *models.Company) error
	Create(ctx context.Context, create *models.Company) error
	Delete(ctx context.Context, id models.UUID) error
}
