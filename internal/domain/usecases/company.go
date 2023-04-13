package usecases

import (
	"context"

	"github.com/018bf/companies/internal/domain/models"
)

// CompanyUseCase - domain layer use case interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/company.go . CompanyUseCase
type CompanyUseCase interface {
	Get(ctx context.Context, id models.UUID) (*models.Company, error)
	List(ctx context.Context, filter *models.CompanyFilter) ([]*models.Company, uint64, error) //deprecated
	Update(ctx context.Context, update *models.CompanyUpdate) (*models.Company, error)
	Create(ctx context.Context, create *models.CompanyCreate) (*models.Company, error)
	Delete(ctx context.Context, id models.UUID) error
}
