package interceptors

import (
	"context"

	"github.com/018bf/companies/internal/domain/interceptors"
	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/internal/domain/usecases"
	"github.com/018bf/companies/pkg/log"
)

type CompanyInterceptor struct {
	companyUseCase usecases.CompanyUseCase
	authUseCase    usecases.AuthUseCase
	logger         log.Logger
}

func NewCompanyInterceptor(
	companyUseCase usecases.CompanyUseCase,
	authUseCase usecases.AuthUseCase,
	logger log.Logger,
) interceptors.CompanyInterceptor {
	return &CompanyInterceptor{
		companyUseCase: companyUseCase,
		authUseCase:    authUseCase,
		logger:         logger,
	}
}

func (i *CompanyInterceptor) Create(
	ctx context.Context,
	create *models.CompanyCreate,
	token *models.Token,
) (*models.Company, error) {
	if err := i.authUseCase.HasPermission(ctx, token, models.PermissionIDCompanyCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, token, models.PermissionIDCompanyCreate, create); err != nil {
		return nil, err
	}
	company, err := i.companyUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (i *CompanyInterceptor) Get(
	ctx context.Context,
	id models.UUID,
	token *models.Token,
) (*models.Company, error) {
	if err := i.authUseCase.HasPermission(ctx, token, models.PermissionIDCompanyDetail); err != nil {
		return nil, err
	}
	company, err := i.companyUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, token, models.PermissionIDCompanyDetail, company); err != nil {
		return nil, err
	}
	return company, nil
}

func (i *CompanyInterceptor) List(
	ctx context.Context,
	filter *models.CompanyFilter,
	token *models.Token,
) ([]*models.Company, uint64, error) {
	if err := i.authUseCase.HasPermission(ctx, token, models.PermissionIDCompanyList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, token, models.PermissionIDCompanyList, filter); err != nil {
		return nil, 0, err
	}
	listCompanies, count, err := i.companyUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return listCompanies, count, nil
}

func (i *CompanyInterceptor) Update(
	ctx context.Context,
	update *models.CompanyUpdate,
	token *models.Token,
) (*models.Company, error) {
	if err := i.authUseCase.HasPermission(ctx, token, models.PermissionIDCompanyUpdate); err != nil {
		return nil, err
	}
	company, err := i.companyUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, token, models.PermissionIDCompanyUpdate, company); err != nil {
		return nil, err
	}
	updated, err := i.companyUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (i *CompanyInterceptor) Delete(
	ctx context.Context,
	id models.UUID,
	token *models.Token,
) error {
	if err := i.authUseCase.HasPermission(ctx, token, models.PermissionIDCompanyDelete); err != nil {
		return err
	}
	company, err := i.companyUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, token, models.PermissionIDCompanyDelete, company); err != nil {
		return err
	}
	if err := i.companyUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
