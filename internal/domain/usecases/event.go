package usecases

import (
	"context"
	"github.com/018bf/companies/internal/domain/models"
)

// EventUseCase - domain layer use case interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/event.go . EventUseCase
type EventUseCase interface {
	CompanyCreated(ctx context.Context, company *models.Company) error
	CompanyUpdated(ctx context.Context, company *models.Company) error
	CompanyDeleted(ctx context.Context, company *models.Company) error
}
