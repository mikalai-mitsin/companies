package interceptor

import (
	"context"

	"github.com/018bf/companies/internal/entity"
	"github.com/018bf/companies/pkg/log"
)

//go:generate mockgen -source=company.go -package=interceptors -destination=company_mock.go

type authService interface {
	HasPermission(ctx context.Context, token *entity.Token, permission entity.PermissionID) error
	HasObjectPermission(
		ctx context.Context,
		token *entity.Token,
		permission entity.PermissionID,
		object any,
	) error
}

type companyService interface {
	Get(ctx context.Context, id entity.UUID) (*entity.Company, error)
	List(ctx context.Context, filter *entity.CompanyFilter) ([]*entity.Company, uint64, error) //deprecated
	Update(ctx context.Context, update *entity.CompanyUpdate) (*entity.Company, error)
	Create(ctx context.Context, create *entity.CompanyCreate) (*entity.Company, error)
	Delete(ctx context.Context, id entity.UUID) error
}

type eventService interface {
	CompanyCreated(ctx context.Context, company *entity.Company) error
	CompanyUpdated(ctx context.Context, company *entity.Company) error
	CompanyDeleted(ctx context.Context, company *entity.Company) error
}

type CompanyInterceptor struct {
	companyService companyService
	authService    authService
	eventService   eventService
	logger         log.Logger
}

func NewCompanyInterceptor(
	companyService companyService,
	authService authService,
	eventService eventService,
	logger log.Logger,
) *CompanyInterceptor {
	return &CompanyInterceptor{
		companyService: companyService,
		authService:    authService,
		eventService:   eventService,
		logger:         logger,
	}
}

func (i *CompanyInterceptor) Create(
	ctx context.Context,
	create *entity.CompanyCreate,
	token *entity.Token,
) (*entity.Company, error) {
	if err := i.authService.HasPermission(ctx, token, entity.PermissionIDCompanyCreate); err != nil {
		return nil, err
	}
	if err := i.authService.HasObjectPermission(ctx, token, entity.PermissionIDCompanyCreate, create); err != nil {
		return nil, err
	}
	company, err := i.companyService.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	if err := i.eventService.CompanyCreated(ctx, company); err != nil {
		i.logger.Error("can't send 'company created' event", log.Context(ctx), log.Error(err))
	}
	return company, nil
}

func (i *CompanyInterceptor) Get(
	ctx context.Context,
	id entity.UUID,
	token *entity.Token,
) (*entity.Company, error) {
	if err := i.authService.HasPermission(ctx, token, entity.PermissionIDCompanyDetail); err != nil {
		return nil, err
	}
	company, err := i.companyService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authService.HasObjectPermission(ctx, token, entity.PermissionIDCompanyDetail, company); err != nil {
		return nil, err
	}
	return company, nil
}

// List
// deprecated
func (i *CompanyInterceptor) List(
	ctx context.Context,
	filter *entity.CompanyFilter,
	token *entity.Token,
) ([]*entity.Company, uint64, error) {
	if err := i.authService.HasPermission(ctx, token, entity.PermissionIDCompanyList); err != nil {
		return nil, 0, err
	}
	if err := i.authService.HasObjectPermission(ctx, token, entity.PermissionIDCompanyList, filter); err != nil {
		return nil, 0, err
	}
	listCompanies, count, err := i.companyService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return listCompanies, count, nil
}

func (i *CompanyInterceptor) Update(
	ctx context.Context,
	update *entity.CompanyUpdate,
	token *entity.Token,
) (*entity.Company, error) {
	if err := i.authService.HasPermission(ctx, token, entity.PermissionIDCompanyUpdate); err != nil {
		return nil, err
	}
	company, err := i.companyService.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authService.HasObjectPermission(ctx, token, entity.PermissionIDCompanyUpdate, company); err != nil {
		return nil, err
	}
	updated, err := i.companyService.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	if err := i.eventService.CompanyUpdated(ctx, company); err != nil {
		i.logger.Error("can't send 'company updated' event", log.Context(ctx), log.Error(err))
	}
	return updated, nil
}

func (i *CompanyInterceptor) Delete(
	ctx context.Context,
	id entity.UUID,
	token *entity.Token,
) error {
	if err := i.authService.HasPermission(ctx, token, entity.PermissionIDCompanyDelete); err != nil {
		return err
	}
	company, err := i.companyService.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authService.HasObjectPermission(ctx, token, entity.PermissionIDCompanyDelete, company); err != nil {
		return err
	}
	if err := i.companyService.Delete(ctx, id); err != nil {
		return err
	}
	if err := i.eventService.CompanyDeleted(ctx, company); err != nil {
		i.logger.Error("can't send 'company deleted' event", log.Context(ctx), log.Error(err))
	}
	return nil
}
