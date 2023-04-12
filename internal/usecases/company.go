package usecases

import (
	"context"

	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/internal/domain/repositories"
	"github.com/018bf/companies/internal/domain/usecases"
	"github.com/018bf/companies/pkg/clock"
	"github.com/018bf/companies/pkg/log"
)

type CompanyUseCase struct {
	companyRepository repositories.CompanyRepository
	clock             clock.Clock
	logger            log.Logger
}

func NewCompanyUseCase(
	companyRepository repositories.CompanyRepository,
	clock clock.Clock,
	logger log.Logger,
) usecases.CompanyUseCase {
	return &CompanyUseCase{companyRepository: companyRepository, clock: clock, logger: logger}
}

func (u *CompanyUseCase) Create(
	ctx context.Context,
	create *models.CompanyCreate,
) (*models.Company, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	company := &models.Company{
		ID:                "",
		UpdatedAt:         now,
		CreatedAt:         now,
		Name:              create.Name,
		Description:       create.Description,
		AmountOfEmployees: create.AmountOfEmployees,
		Registered:        create.Registered,
		Type:              create.Type,
	}
	if err := u.companyRepository.Create(ctx, company); err != nil {
		return nil, err
	}
	return company, nil
}
func (u *CompanyUseCase) Get(ctx context.Context, id models.UUID) (*models.Company, error) {
	if err := id.Validate(); err != nil {
		return nil, err
	}
	company, err := u.companyRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (u *CompanyUseCase) List(
	ctx context.Context,
	filter *models.CompanyFilter,
) ([]*models.Company, uint64, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, err
	}
	company, err := u.companyRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.companyRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return company, count, nil
}

func (u *CompanyUseCase) Update(
	ctx context.Context,
	update *models.CompanyUpdate,
) (*models.Company, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	company, err := u.companyRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if update.Name != nil {
		company.Name = *update.Name
	}
	if update.Description != nil {
		company.Description = *update.Description
	}
	if update.AmountOfEmployees != nil {
		company.AmountOfEmployees = *update.AmountOfEmployees
	}
	if update.Registered != nil {
		company.Registered = *update.Registered
	}
	if update.Type != nil {
		company.Type = *update.Type
	}
	company.UpdatedAt = u.clock.Now().UTC()
	if err := u.companyRepository.Update(ctx, company); err != nil {
		return nil, err
	}
	return company, nil
}
func (u *CompanyUseCase) Delete(ctx context.Context, id models.UUID) error {
	if err := u.companyRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
