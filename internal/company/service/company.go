package service

import (
	"context"

	"github.com/018bf/companies/internal/entity"
	"github.com/018bf/companies/pkg/clock"
	"github.com/018bf/companies/pkg/log"
)

//go:generate mockgen -source=company.go -package=usecases -destination=company_mock.go

type companyRepository interface {
	Get(ctx context.Context, id entity.UUID) (*entity.Company, error)
	List(ctx context.Context, filter *entity.CompanyFilter) ([]*entity.Company, error) // deprecated
	Count(ctx context.Context, filter *entity.CompanyFilter) (uint64, error)           // deprecated
	Update(ctx context.Context, update *entity.Company) error
	Create(ctx context.Context, create *entity.Company) error
	Delete(ctx context.Context, id entity.UUID) error
}

type CompanyService struct {
	companyRepository companyRepository
	clock             clock.Clock
	logger            log.Logger
}

func NewCompanyService(
	companyRepository companyRepository,
	clock clock.Clock,
	logger log.Logger,
) *CompanyService {
	return &CompanyService{companyRepository: companyRepository, clock: clock, logger: logger}
}

func (u *CompanyService) Create(
	ctx context.Context,
	create *entity.CompanyCreate,
) (*entity.Company, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	company := &entity.Company{
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
func (u *CompanyService) Get(ctx context.Context, id entity.UUID) (*entity.Company, error) {
	if err := id.Validate(); err != nil {
		return nil, err
	}
	company, err := u.companyRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return company, nil
}

// List
// deprecated
func (u *CompanyService) List(
	ctx context.Context,
	filter *entity.CompanyFilter,
) ([]*entity.Company, uint64, error) {
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

func (u *CompanyService) Update(
	ctx context.Context,
	update *entity.CompanyUpdate,
) (*entity.Company, error) {
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
func (u *CompanyService) Delete(ctx context.Context, id entity.UUID) error {
	if err := u.companyRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
