package integration

import (
	"context"
	mock_models "github.com/018bf/companies/internal/domain/models/mock"
	companiespb "github.com/018bf/companies/pkg/companiespb/v1"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"syreclabs.com/go/faker"
	"testing"
)

func TestRegionCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	client := companiespb.NewCompanyServiceClient(conn)
	ctx := context.Background()
	create := mock_models.NewCompanyCreate(t)
	type fields struct {
		client companiespb.CompanyServiceClient
	}
	type args struct {
		ctx   context.Context
		input *companiespb.CompanyCreate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *companiespb.Company
		wantErr *status.Status
	}{
		{
			name: "ok",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: getAuthContext(ctx, accessToken),
				input: &companiespb.CompanyCreate{
					Name:              create.Name,
					Description:       create.Description,
					AmountOfEmployees: int32(create.AmountOfEmployees),
					Registered:        create.Registered,
					Type:              companiespb.CompanyType(create.Type),
				},
			},
			want: &companiespb.Company{
				Id:                "",
				UpdatedAt:         nil,
				CreatedAt:         nil,
				Name:              create.Name,
				Description:       create.Description,
				AmountOfEmployees: int32(create.AmountOfEmployees),
				Registered:        create.Registered,
				Type:              companiespb.CompanyType(create.Type),
			},
			wantErr: nil,
		},
		{
			name: "permission error",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: context.Background(),
				input: &companiespb.CompanyCreate{
					Name: company.Name,
				},
			},
			want:    nil,
			wantErr: status.New(codes.PermissionDenied, "Permission denied."),
		},
		{
			name: "invalid form",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: getAuthContext(ctx, accessToken),
				input: &companiespb.CompanyCreate{
					Description:       company.Description,
					AmountOfEmployees: int32(company.AmountOfEmployees),
					Registered:        company.Registered,
					Type:              companiespb.CompanyType(company.Type),
				},
			},
			want:    nil,
			wantErr: status.New(codes.InvalidArgument, "The form sent is not valid, please correct the errors below."),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.client.Create(tt.args.ctx, tt.args.input)
			s, _ := status.FromError(err)
			if !statusEqual(s, tt.wantErr) {
				t.Errorf("Create() err = %v, wantErr %v", err, tt.wantErr)
			}
			if got.GetName() != tt.want.GetName() {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
			if got.GetDescription() != tt.want.GetDescription() {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
			if got.GetAmountOfEmployees() != tt.want.GetAmountOfEmployees() {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
			if got.GetAmountOfEmployees() != tt.want.GetAmountOfEmployees() {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
			if got.GetRegistered() != tt.want.GetRegistered() {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
			if got.GetType() != tt.want.GetType() {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegionUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	client := companiespb.NewCompanyServiceClient(conn)
	ctx := context.Background()
	update := mock_models.NewCompanyUpdate(t)
	type fields struct {
		client companiespb.CompanyServiceClient
	}
	type args struct {
		ctx   context.Context
		input *companiespb.CompanyUpdate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *companiespb.Company
		wantErr *status.Status
	}{
		{
			name: "ok",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: getAuthContext(ctx, accessToken),
				input: &companiespb.CompanyUpdate{
					Id:                string(company.ID),
					Name:              wrapperspb.String(*update.Name),
					Description:       wrapperspb.String(*update.Description),
					AmountOfEmployees: wrapperspb.Int32(int32(*update.AmountOfEmployees)),
					Registered:        wrapperspb.Bool(*update.Registered),
					Type:              companiespb.CompanyType(*update.Type),
				},
			},
			want: &companiespb.Company{
				Id:                string(company.ID),
				UpdatedAt:         timestamppb.New(company.UpdatedAt),
				CreatedAt:         timestamppb.New(company.CreatedAt),
				Name:              *update.Name,
				Description:       *update.Description,
				AmountOfEmployees: int32(*update.AmountOfEmployees),
				Registered:        *update.Registered,
				Type:              companiespb.CompanyType(*update.Type),
			},
			wantErr: nil,
		},
		{
			name: "not found",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: getAuthContext(ctx, accessToken),
				input: &companiespb.CompanyUpdate{
					Id:                "d77a0a83-83f9-436b-a462-532b70b44f9b",
					Name:              wrapperspb.String(*update.Name),
					Description:       wrapperspb.String(*update.Description),
					AmountOfEmployees: wrapperspb.Int32(int32(*update.AmountOfEmployees)),
					Registered:        wrapperspb.Bool(*update.Registered),
					Type:              companiespb.CompanyType(*update.Type),
				},
			},
			want:    nil,
			wantErr: status.New(codes.NotFound, "Region not found."),
		},
		{
			name: "invalid id",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: getAuthContext(ctx, accessToken),
				input: &companiespb.CompanyUpdate{
					Id:                faker.Lorem().String(),
					Name:              wrapperspb.String(*update.Name),
					Description:       wrapperspb.String(*update.Description),
					AmountOfEmployees: wrapperspb.Int32(int32(*update.AmountOfEmployees)),
					Registered:        wrapperspb.Bool(*update.Registered),
					Type:              companiespb.CompanyType(*update.Type),
				},
			},
			want:    nil,
			wantErr: status.New(codes.InvalidArgument, "The form sent is not valid, please correct the errors below."),
		},
		{
			name: "permission error",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: context.Background(),
				input: &companiespb.CompanyUpdate{
					Id:                "d77a0a83-83f9-436b-a462-532b70b44f9b",
					Name:              wrapperspb.String(*update.Name),
					Description:       wrapperspb.String(*update.Description),
					AmountOfEmployees: wrapperspb.Int32(int32(*update.AmountOfEmployees)),
					Registered:        wrapperspb.Bool(*update.Registered),
					Type:              companiespb.CompanyType(*update.Type),
				},
			},
			want:    nil,
			wantErr: status.New(codes.PermissionDenied, "Permission denied."),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.client.Update(tt.args.ctx, tt.args.input)
			s, _ := status.FromError(err)
			if !statusEqual(s, tt.wantErr) {
				t.Errorf("Update() err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != nil && tt.want != nil {
				tt.want.UpdatedAt = got.UpdatedAt
			}
			if got.GetId() != tt.want.GetId() {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
			if got.GetName() != tt.want.GetName() {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
			if got.GetDescription() != tt.want.GetDescription() {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
			if got.GetAmountOfEmployees() != tt.want.GetAmountOfEmployees() {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
			if got.GetAmountOfEmployees() != tt.want.GetAmountOfEmployees() {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
			if got.GetRegistered() != tt.want.GetRegistered() {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
			if got.GetType() != tt.want.GetType() {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegionDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	client := companiespb.NewCompanyServiceClient(conn)
	ctx := context.Background()
	type fields struct {
		client companiespb.CompanyServiceClient
	}
	type args struct {
		ctx   context.Context
		input *companiespb.CompanyDelete
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr *status.Status
	}{
		{
			name: "ok",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: getAuthContext(ctx, accessToken),
				input: &companiespb.CompanyDelete{
					Id: string(company.ID),
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "not found",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: getAuthContext(ctx, accessToken),
				input: &companiespb.CompanyDelete{
					Id: "d77a0a83-83f9-436b-a462-532b70b44f9b",
				},
			},
			want:    nil,
			wantErr: status.New(codes.NotFound, "Entity not found."),
		},
		{
			name: "invalid argument",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: getAuthContext(ctx, accessToken),
				input: &companiespb.CompanyDelete{
					Id: faker.Lorem().String(),
				},
			},
			want: nil,
			wantErr: status.
				New(codes.InvalidArgument, "The form sent is not valid, please correct the errors below."),
		},
		{
			name: "permission error",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: context.Background(),
				input: &companiespb.CompanyDelete{
					Id: string(company.ID),
				},
			},
			want:    nil,
			wantErr: status.New(codes.PermissionDenied, "Permission denied."),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.client.Delete(tt.args.ctx, tt.args.input)
			s, _ := status.FromError(err)
			if !statusEqual(s, tt.wantErr) {
				t.Errorf("Delete() err = %v, wantErr %v", err, tt.wantErr)
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegionList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	client := companiespb.NewCompanyServiceClient(conn)
	ctx := context.Background()
	type fields struct {
		client companiespb.CompanyServiceClient
	}
	type args struct {
		ctx   context.Context
		input *companiespb.CompanyFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr *status.Status
	}{
		{
			name: "ok",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: getAuthContext(ctx, accessToken),
				input: &companiespb.CompanyFilter{
					PageNumber: wrapperspb.UInt64(1),
					PageSize:   wrapperspb.UInt64(5),
					Search:     nil,
					Registered: wrapperspb.Bool(true),
					OrderBy:    nil,
					Ids:        nil,
					Types: []companiespb.CompanyType{
						companiespb.CompanyType_COMPANY_TYPE_COOPERATIVE,
					},
				},
			},
			wantErr: status.New(codes.PermissionDenied, "Permission denied."),
		},
		{
			name: "invalid argument",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: getAuthContext(ctx, accessToken),
				input: &companiespb.CompanyFilter{
					PageNumber: wrapperspb.UInt64(1),
					PageSize:   wrapperspb.UInt64(5),
					Search:     wrapperspb.String("search"),
				},
			},
			wantErr: status.New(codes.PermissionDenied, "Permission denied."),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.client.List(tt.args.ctx, tt.args.input) // nolint: staticcheck
			s, _ := status.FromError(err)
			if !statusEqual(s, tt.wantErr) {
				t.Errorf("List() err = %v, wantErr %v", err, tt.wantErr)
			}
			for _, company := range got.GetItems() {
				if !slices.Contains(tt.args.input.GetTypes(), company.GetType()) {
					t.Errorf("List() got type = %v, want type %v", company.GetType(), tt.args.input.GetTypes())
				}
				if company.GetRegistered() != tt.args.input.GetRegistered().GetValue() {
					t.Errorf("Update() got registered = %v, want registered %v", got, tt.args.input.GetRegistered().GetValue())
				}
			}
		})
	}
}
