package grpc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/018bf/companies/internal/domain/errs"
	"github.com/018bf/companies/internal/domain/interceptors"
	mock_interceptors "github.com/018bf/companies/internal/domain/interceptors/mock"
	"github.com/018bf/companies/internal/domain/models"
	mock_models "github.com/018bf/companies/internal/domain/models/mock"
	companiespb "github.com/018bf/companies/pkg/companiespb/v1"
	"github.com/018bf/companies/pkg/log"
	mock_log "github.com/018bf/companies/pkg/log/mock"
	"github.com/018bf/companies/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewCompanyServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	companyInterceptor := mock_interceptors.NewMockCompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		companyInterceptor interceptors.CompanyInterceptor
		logger             log.Logger
	}
	tests := []struct {
		name string
		args args
		want companiespb.CompanyServiceServer
	}{
		{
			name: "ok",
			args: args{
				companyInterceptor: companyInterceptor,
				logger:             logger,
			},
			want: &CompanyServiceServer{
				companyInterceptor: companyInterceptor,
				logger:             logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCompanyServiceServer(tt.args.companyInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewCompanyServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	companyInterceptor := mock_interceptors.NewMockCompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := utils.Pointer(mock_models.NewToken(t))
	ctx = context.WithValue(ctx, TokenKey, user)
	// create := mock_models.NewCompanyCreate(t)
	company := mock_models.NewCompany(t)
	type fields struct {
		UnimplementedCompanyServiceServer companiespb.UnimplementedCompanyServiceServer
		companyInterceptor                interceptors.CompanyInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *companiespb.CompanyCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *companiespb.Company
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				companyInterceptor.
					EXPECT().
					Create(ctx, gomock.Any(), user).
					Return(company, nil)
			},
			fields: fields{
				UnimplementedCompanyServiceServer: companiespb.UnimplementedCompanyServiceServer{},
				companyInterceptor:                companyInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx:   ctx,
				input: &companiespb.CompanyCreate{},
			},
			want:    decodeCompany(company),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				companyInterceptor.
					EXPECT().
					Create(ctx, gomock.Any(), user).
					Return(nil, errs.NewUnexpectedBehaviorError("interceptor error")).
					Times(1)
			},
			fields: fields{
				UnimplementedCompanyServiceServer: companiespb.UnimplementedCompanyServiceServer{},
				companyInterceptor:                companyInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx:   ctx,
				input: &companiespb.CompanyCreate{},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("interceptor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := CompanyServiceServer{
				UnimplementedCompanyServiceServer: tt.fields.UnimplementedCompanyServiceServer,
				companyInterceptor:                tt.fields.companyInterceptor,
				logger:                            tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	companyInterceptor := mock_interceptors.NewMockCompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := utils.Pointer(mock_models.NewToken(t))
	ctx = context.WithValue(ctx, TokenKey, user)
	id := uuid.NewString()
	type fields struct {
		UnimplementedCompanyServiceServer companiespb.UnimplementedCompanyServiceServer
		companyInterceptor                interceptors.CompanyInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *companiespb.CompanyDelete
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				companyInterceptor.EXPECT().Delete(ctx, models.UUID(id), user).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedCompanyServiceServer: companiespb.UnimplementedCompanyServiceServer{},
				companyInterceptor:                companyInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &companiespb.CompanyDelete{
					Id: id,
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				companyInterceptor.EXPECT().Delete(ctx, models.UUID(id), user).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedCompanyServiceServer: companiespb.UnimplementedCompanyServiceServer{},
				companyInterceptor:                companyInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &companiespb.CompanyDelete{
					Id: id,
				},
			},
			want: nil,
			wantErr: decodeError(&errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"details": "i error",
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := CompanyServiceServer{
				UnimplementedCompanyServiceServer: tt.fields.UnimplementedCompanyServiceServer,
				companyInterceptor:                tt.fields.companyInterceptor,
				logger:                            tt.fields.logger,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	companyInterceptor := mock_interceptors.NewMockCompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := utils.Pointer(mock_models.NewToken(t))
	ctx = context.WithValue(ctx, TokenKey, user)
	company := mock_models.NewCompany(t)
	type fields struct {
		UnimplementedCompanyServiceServer companiespb.UnimplementedCompanyServiceServer
		companyInterceptor                interceptors.CompanyInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *companiespb.CompanyGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *companiespb.Company
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				companyInterceptor.EXPECT().Get(ctx, company.ID, user).Return(company, nil).Times(1)
			},
			fields: fields{
				UnimplementedCompanyServiceServer: companiespb.UnimplementedCompanyServiceServer{},
				companyInterceptor:                companyInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &companiespb.CompanyGet{
					Id: string(company.ID),
				},
			},
			want:    decodeCompany(company),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				companyInterceptor.EXPECT().Get(ctx, company.ID, user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedCompanyServiceServer: companiespb.UnimplementedCompanyServiceServer{},
				companyInterceptor:                companyInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &companiespb.CompanyGet{
					Id: string(company.ID),
				},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := CompanyServiceServer{
				UnimplementedCompanyServiceServer: tt.fields.UnimplementedCompanyServiceServer,
				companyInterceptor:                tt.fields.companyInterceptor,
				logger:                            tt.fields.logger,
			}
			got, err := s.Get(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	companyInterceptor := mock_interceptors.NewMockCompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := utils.Pointer(mock_models.NewToken(t))
	ctx = context.WithValue(ctx, TokenKey, user)
	filter := mock_models.NewCompanyFilter(t)
	filter.Types = nil
	var ids []models.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &companiespb.ListCompany{
		Items: make([]*companiespb.Company, 0, int(count)),
		Count: count,
	}
	listCompanies := make([]*models.Company, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := mock_models.NewCompany(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, string(a.ID))
		listCompanies = append(listCompanies, a)
		response.Items = append(response.Items, decodeCompany(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedCompanyServiceServer companiespb.UnimplementedCompanyServiceServer
		companyInterceptor                interceptors.CompanyInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *companiespb.CompanyFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *companiespb.ListCompany
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				companyInterceptor.EXPECT().
					List(ctx, filter, user).
					Return(listCompanies, count, nil)
			},
			fields: fields{
				UnimplementedCompanyServiceServer: companiespb.UnimplementedCompanyServiceServer{},
				companyInterceptor:                companyInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &companiespb.CompanyFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					Search:     wrapperspb.String(*filter.Search),
					Registered: wrapperspb.Bool(*filter.Registered),
					OrderBy:    filter.OrderBy,
					Ids:        stringIDs,
					Types:      nil,
				},
			},
			want:    response,
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				companyInterceptor.
					EXPECT().
					List(ctx, filter, user).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedCompanyServiceServer: companiespb.UnimplementedCompanyServiceServer{},
				companyInterceptor:                companyInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &companiespb.CompanyFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					Search:     wrapperspb.String(*filter.Search),
					Registered: wrapperspb.Bool(*filter.Registered),
					OrderBy:    filter.OrderBy,
					Ids:        stringIDs,
					Types:      nil,
				},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := CompanyServiceServer{
				UnimplementedCompanyServiceServer: tt.fields.UnimplementedCompanyServiceServer,
				companyInterceptor:                tt.fields.companyInterceptor,
				logger:                            tt.fields.logger,
			}
			got, err := s.List(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	companyInterceptor := mock_interceptors.NewMockCompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := utils.Pointer(mock_models.NewToken(t))
	ctx = context.WithValue(ctx, TokenKey, user)
	company := mock_models.NewCompany(t)
	update := mock_models.NewCompanyUpdate(t)
	type fields struct {
		UnimplementedCompanyServiceServer companiespb.UnimplementedCompanyServiceServer
		companyInterceptor                interceptors.CompanyInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *companiespb.CompanyUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *companiespb.Company
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				companyInterceptor.EXPECT().
					Update(ctx, gomock.Any(), user).
					Return(company, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedCompanyServiceServer: companiespb.UnimplementedCompanyServiceServer{},
				companyInterceptor:                companyInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeCompanyUpdate(update),
			},
			want:    decodeCompany(company),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				companyInterceptor.EXPECT().Update(ctx, gomock.Any(), user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedCompanyServiceServer: companiespb.UnimplementedCompanyServiceServer{},
				companyInterceptor:                companyInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeCompanyUpdate(update),
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := CompanyServiceServer{
				UnimplementedCompanyServiceServer: tt.fields.UnimplementedCompanyServiceServer,
				companyInterceptor:                tt.fields.companyInterceptor,
				logger:                            tt.fields.logger,
			}
			got, err := s.Update(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeCompany(t *testing.T) {
	company := mock_models.NewCompany(t)
	result := &companiespb.Company{
		Id:                string(company.ID),
		UpdatedAt:         timestamppb.New(company.UpdatedAt),
		CreatedAt:         timestamppb.New(company.CreatedAt),
		Name:              company.Name,
		Description:       company.Description,
		AmountOfEmployees: int32(company.AmountOfEmployees),
		Registered:        company.Registered,
		Type:              decodeCompanyType(company.Type),
	}
	type args struct {
		company *models.Company
	}
	tests := []struct {
		name string
		args args
		want *companiespb.Company
	}{
		{
			name: "ok",
			args: args{
				company: company,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeCompany(tt.args.company); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeCompany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeCompanyFilter(t *testing.T) {
	id := models.UUID(uuid.NewString())
	type args struct {
		input *companiespb.CompanyFilter
	}
	tests := []struct {
		name string
		args args
		want *models.CompanyFilter
	}{
		{
			name: "ok",
			args: args{
				input: &companiespb.CompanyFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					Search:     wrapperspb.String("my name is"),
					OrderBy:    []string{"created_at", "id"},
					Ids:        []string{string(id)},
				},
			},
			want: &models.CompanyFilter{
				PageSize:   utils.Pointer(uint64(5)),
				PageNumber: utils.Pointer(uint64(2)),
				OrderBy:    []string{"created_at", "id"},
				Search:     utils.Pointer("my name is"),
				IDs:        []models.UUID{id},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeCompanyFilter(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeUserFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
