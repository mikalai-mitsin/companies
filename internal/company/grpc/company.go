package grpc

import (
	"context"
	grpc2 "github.com/018bf/companies/internal/interfaces/grpc"

	"github.com/018bf/companies/internal/entity"
	companiespb "github.com/018bf/companies/pkg/companiespb/v1"
	"github.com/018bf/companies/pkg/log"
	"github.com/018bf/companies/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

//go:generate mockgen -source=company.go -package=grpc -destination=company_mock.go

type companyInterceptor interface {
	Get(ctx context.Context, id entity.UUID, token *entity.Token) (*entity.Company, error)
	List(
		ctx context.Context,
		filter *entity.CompanyFilter,
		token *entity.Token,
	) ([]*entity.Company, uint64, error) // deprecated
	Update(
		ctx context.Context,
		update *entity.CompanyUpdate,
		token *entity.Token,
	) (*entity.Company, error)
	Create(
		ctx context.Context,
		create *entity.CompanyCreate,
		token *entity.Token,
	) (*entity.Company, error)
	Delete(ctx context.Context, id entity.UUID, token *entity.Token) error
}

type CompanyServiceServer struct {
	companiespb.UnimplementedCompanyServiceServer
	companyInterceptor companyInterceptor
	logger             log.Logger
}

func NewCompanyServiceServer(
	companyInterceptor companyInterceptor,
	logger log.Logger,
) *CompanyServiceServer {
	return &CompanyServiceServer{companyInterceptor: companyInterceptor, logger: logger}
}

func (s *CompanyServiceServer) Create(
	ctx context.Context,
	input *companiespb.CompanyCreate,
) (*companiespb.Company, error) {
	company, err := s.companyInterceptor.Create(
		ctx,
		encodeCompanyCreate(input),
		ctx.Value(grpc2.TokenKey).(*entity.Token),
	)
	if err != nil {
		return nil, grpc2.DecodeError(err)
	}
	return decodeCompany(company), nil
}

func (s *CompanyServiceServer) Get(
	ctx context.Context,
	input *companiespb.CompanyGet,
) (*companiespb.Company, error) {
	company, err := s.companyInterceptor.Get(
		ctx,
		entity.UUID(input.GetId()),
		ctx.Value(grpc2.TokenKey).(*entity.Token),
	)
	if err != nil {
		return nil, grpc2.DecodeError(err)
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
		ctx.Value(grpc2.TokenKey).(*entity.Token),
	)
	if err != nil {
		return nil, grpc2.DecodeError(err)
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
		ctx.Value(grpc2.TokenKey).(*entity.Token),
	)
	if err != nil {
		return nil, grpc2.DecodeError(err)
	}
	return decodeCompany(company), nil
}

func (s *CompanyServiceServer) Delete(
	ctx context.Context,
	input *companiespb.CompanyDelete,
) (*emptypb.Empty, error) {
	if err := s.companyInterceptor.Delete(
		ctx,
		entity.UUID(input.GetId()),
		ctx.Value(grpc2.TokenKey).(*entity.Token),
	); err != nil {
		return nil, grpc2.DecodeError(err)
	}
	return &emptypb.Empty{}, nil
}
func encodeCompanyCreate(input *companiespb.CompanyCreate) *entity.CompanyCreate {
	create := &entity.CompanyCreate{
		Name:              input.GetName(),
		Description:       input.GetDescription(),
		AmountOfEmployees: int(input.GetAmountOfEmployees()),
		Registered:        input.GetRegistered(),
		Type:              encodeCompanyType(input.GetType()),
	}
	return create
}
func encodeCompanyType(companyType companiespb.CompanyType) entity.CompanyType {
	switch companyType {
	case companiespb.CompanyType_COMPANY_TYPE_CORPORATIONS:
		return entity.CompanyTypeCorporations
	case companiespb.CompanyType_COMPANY_TYPE_NON_PROFIT:
		return entity.CompanyTypeNonProfit
	case companiespb.CompanyType_COMPANY_TYPE_COOPERATIVE:
		return entity.CompanyTypeCooperative
	case companiespb.CompanyType_COMPANY_TYPE_SOLE_PROPRIETORSHIP:
		return entity.CompanyTypeSoleProprietorship
	default:
		return 0
	}
}
func encodeCompanyFilter(input *companiespb.CompanyFilter) *entity.CompanyFilter {
	filter := &entity.CompanyFilter{
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
		filter.IDs = append(filter.IDs, entity.UUID(id))
	}
	if input.GetSearch() != nil {
		filter.Search = utils.Pointer(input.GetSearch().GetValue())
	}
	if len(input.GetTypes()) > 0 {
		filter.Types = make([]entity.CompanyType, len(input.GetTypes()))
		for i, companyType := range input.GetTypes() {
			filter.Types[i] = encodeCompanyType(companyType)
		}
	}
	if input.GetRegistered() != nil {
		filter.Registered = utils.Pointer(input.GetRegistered().GetValue())
	}
	return filter
}
func encodeCompanyUpdate(input *companiespb.CompanyUpdate) *entity.CompanyUpdate {
	update := &entity.CompanyUpdate{ID: entity.UUID(input.GetId())}
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
func decodeCompany(company *entity.Company) *companiespb.Company {
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
func decodeCompanyType(companyType entity.CompanyType) companiespb.CompanyType {
	switch companyType {
	case entity.CompanyTypeCorporations:
		return companiespb.CompanyType_COMPANY_TYPE_CORPORATIONS
	case entity.CompanyTypeNonProfit:
		return companiespb.CompanyType_COMPANY_TYPE_NON_PROFIT
	case entity.CompanyTypeCooperative:
		return companiespb.CompanyType_COMPANY_TYPE_COOPERATIVE
	case entity.CompanyTypeSoleProprietorship:
		return companiespb.CompanyType_COMPANY_TYPE_SOLE_PROPRIETORSHIP
	default:
		return 0
	}
}
func decodeListCompany(listCompanies []*entity.Company, count uint64) *companiespb.ListCompany {
	response := &companiespb.ListCompany{
		Items: make([]*companiespb.Company, 0, len(listCompanies)),
		Count: count,
	}
	for _, company := range listCompanies {
		response.Items = append(response.Items, decodeCompany(company))
	}
	return response
}
func decodeCompanyUpdate(update *entity.CompanyUpdate) *companiespb.CompanyUpdate {
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
