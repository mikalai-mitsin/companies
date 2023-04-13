package grpc

import (
	"context"

	"github.com/018bf/companies/internal/domain/interceptors"
	"github.com/018bf/companies/internal/domain/models"
	companiespb "github.com/018bf/companies/pkg/companiespb/v1"
	"github.com/018bf/companies/pkg/log"
	"github.com/018bf/companies/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CompanyServiceServer struct {
	companiespb.UnimplementedCompanyServiceServer
	companyInterceptor interceptors.CompanyInterceptor
	logger             log.Logger
}

func NewCompanyServiceServer(
	companyInterceptor interceptors.CompanyInterceptor,
	logger log.Logger,
) companiespb.CompanyServiceServer {
	return &CompanyServiceServer{companyInterceptor: companyInterceptor, logger: logger}
}

func (s *CompanyServiceServer) Create(
	ctx context.Context,
	input *companiespb.CompanyCreate,
) (*companiespb.Company, error) {
	company, err := s.companyInterceptor.Create(
		ctx,
		encodeCompanyCreate(input),
		ctx.Value(TokenKey).(*models.Token),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeCompany(company), nil
}

func (s *CompanyServiceServer) Get(
	ctx context.Context,
	input *companiespb.CompanyGet,
) (*companiespb.Company, error) {
	company, err := s.companyInterceptor.Get(
		ctx,
		models.UUID(input.GetId()),
		ctx.Value(TokenKey).(*models.Token),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeCompany(company), nil
}

// List
// deprecated
func (s *CompanyServiceServer) List(
	ctx context.Context,
	filter *companiespb.CompanyFilter,
) (*companiespb.ListCompany, error) {
	listCompanies, count, err := s.companyInterceptor.List(
		ctx,
		encodeCompanyFilter(filter),
		ctx.Value(TokenKey).(*models.Token),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeListCompany(listCompanies, count), nil
}

func (s *CompanyServiceServer) Update(
	ctx context.Context,
	input *companiespb.CompanyUpdate,
) (*companiespb.Company, error) {
	company, err := s.companyInterceptor.Update(
		ctx,
		encodeCompanyUpdate(input),
		ctx.Value(TokenKey).(*models.Token),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeCompany(company), nil
}

func (s *CompanyServiceServer) Delete(
	ctx context.Context,
	input *companiespb.CompanyDelete,
) (*emptypb.Empty, error) {
	if err := s.companyInterceptor.Delete(
		ctx,
		models.UUID(input.GetId()),
		ctx.Value(TokenKey).(*models.Token),
	); err != nil {
		return nil, decodeError(err)
	}
	return &emptypb.Empty{}, nil
}
func encodeCompanyCreate(input *companiespb.CompanyCreate) *models.CompanyCreate {
	create := &models.CompanyCreate{
		Name:              input.GetName(),
		Description:       input.GetDescription(),
		AmountOfEmployees: int(input.GetAmountOfEmployees()),
		Registered:        input.GetRegistered(),
		Type:              encodeCompanyType(input.GetType()),
	}
	return create
}
func encodeCompanyType(companyType companiespb.CompanyType) models.CompanyType {
	switch companyType {
	case companiespb.CompanyType_COMPANY_TYPE_CORPORATIONS:
		return models.CompanyTypeCorporations
	case companiespb.CompanyType_COMPANY_TYPE_NON_PROFIT:
		return models.CompanyTypeNonProfit
	case companiespb.CompanyType_COMPANY_TYPE_COOPERATIVE:
		return models.CompanyTypeCooperative
	case companiespb.CompanyType_COMPANY_TYPE_SOLE_PROPRIETORSHIP:
		return models.CompanyTypeSoleProprietorship
	default:
		return 0
	}
}
func encodeCompanyFilter(input *companiespb.CompanyFilter) *models.CompanyFilter {
	filter := &models.CompanyFilter{
		IDs:        nil,
		PageSize:   nil,
		PageNumber: nil,
		OrderBy:    input.GetOrderBy(),
		Search:     nil,
	}
	if input.GetPageSize() != nil {
		filter.PageSize = utils.Pointer(input.GetPageSize().GetValue())
	}
	if input.GetPageNumber() != nil {
		filter.PageNumber = utils.Pointer(input.GetPageNumber().GetValue())
	}
	for _, id := range input.GetIds() {
		filter.IDs = append(filter.IDs, models.UUID(id))
	}
	if input.GetSearch() != nil {
		filter.Search = utils.Pointer(input.GetSearch().GetValue())
	}
	if len(input.GetTypes()) > 0 {
		filter.Types = make([]models.CompanyType, len(input.GetTypes()))
		for i, companyType := range input.GetTypes() {
			filter.Types[i] = encodeCompanyType(companyType)
		}
	}
	if input.GetRegistered() != nil {
		filter.Registered = utils.Pointer(input.GetRegistered().GetValue())
	}
	return filter
}
func encodeCompanyUpdate(input *companiespb.CompanyUpdate) *models.CompanyUpdate {
	update := &models.CompanyUpdate{ID: models.UUID(input.GetId())}
	if input.GetName() != nil {
		update.Name = utils.Pointer(input.GetName().GetValue())
	}
	if input.GetDescription() != nil {
		update.Description = utils.Pointer(input.GetDescription().GetValue())
	}
	if input.GetAmountOfEmployees() != nil {
		update.AmountOfEmployees = utils.Pointer(int(input.GetAmountOfEmployees().GetValue()))
	}
	if input.GetRegistered() != nil {
		update.Registered = utils.Pointer(input.GetRegistered().GetValue())
	}
	if input.GetType() != companiespb.CompanyType_COMPANY_TYPE_UNKNOWN {
		update.Type = utils.Pointer(encodeCompanyType(input.GetType()))
	}
	return update
}
func decodeCompany(company *models.Company) *companiespb.Company {
	response := &companiespb.Company{
		Id:                string(company.ID),
		UpdatedAt:         timestamppb.New(company.UpdatedAt),
		CreatedAt:         timestamppb.New(company.CreatedAt),
		Name:              company.Name,
		Description:       company.Description,
		AmountOfEmployees: int32(company.AmountOfEmployees),
		Registered:        company.Registered,
		Type:              decodeCompanyType(company.Type),
	}
	return response
}
func decodeCompanyType(companyType models.CompanyType) companiespb.CompanyType {
	switch companyType {
	case models.CompanyTypeCorporations:
		return companiespb.CompanyType_COMPANY_TYPE_CORPORATIONS
	case models.CompanyTypeNonProfit:
		return companiespb.CompanyType_COMPANY_TYPE_NON_PROFIT
	case models.CompanyTypeCooperative:
		return companiespb.CompanyType_COMPANY_TYPE_COOPERATIVE
	case models.CompanyTypeSoleProprietorship:
		return companiespb.CompanyType_COMPANY_TYPE_SOLE_PROPRIETORSHIP
	default:
		return 0
	}
}
func decodeListCompany(listCompanies []*models.Company, count uint64) *companiespb.ListCompany {
	response := &companiespb.ListCompany{
		Items: make([]*companiespb.Company, 0, len(listCompanies)),
		Count: count,
	}
	for _, company := range listCompanies {
		response.Items = append(response.Items, decodeCompany(company))
	}
	return response
}
func decodeCompanyUpdate(update *models.CompanyUpdate) *companiespb.CompanyUpdate {
	result := &companiespb.CompanyUpdate{
		Id:                string(update.ID),
		Name:              wrapperspb.String(*update.Name),
		Description:       wrapperspb.String(*update.Description),
		AmountOfEmployees: wrapperspb.Int32(int32(*update.AmountOfEmployees)),
		Registered:        wrapperspb.Bool(*update.Registered),
		Type:              decodeCompanyType(*update.Type),
	}
	return result
}
