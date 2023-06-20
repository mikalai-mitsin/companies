package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/018bf/companies/internal/entity"
	mock_models "github.com/018bf/companies/internal/entity/mock"
	"github.com/018bf/companies/internal/errs"
	"github.com/018bf/companies/pkg/clock"
	mock_clock "github.com/018bf/companies/pkg/clock/mock"
	"github.com/018bf/companies/pkg/log"
	mock_log "github.com/018bf/companies/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"github.com/jaswdr/faker"
)

func TestNewCompanyService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyRepository := NewMockcompanyRepository(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		companyRepository companyRepository
		clock             clock.Clock
		logger            log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *CompanyService
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				companyRepository: mockCompanyRepository,
				clock:             clockMock,
				logger:            logger,
			},
			want: &CompanyService{
				companyRepository: mockCompanyRepository,
				clock:             clockMock,
				logger:            logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewCompanyService(tt.args.companyRepository, tt.args.clock, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewCompanyService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyRepository := NewMockcompanyRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	type fields struct {
		companyRepository companyRepository
		logger            log.Logger
	}
	type args struct {
		ctx context.Context
		id  entity.UUID
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *entity.Company
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockCompanyRepository.EXPECT().Get(ctx, company.ID).Return(company, nil)
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				id:  company.ID,
			},
			want:    company,
			wantErr: nil,
		},
		{
			name:  "invalid id",
			setup: func() {},
			fields: fields{
				companyRepository: mockCompanyRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				id:  "company.ID",
			},
			want: nil,
			wantErr: &errs.Error{
				Code:    3,
				Message: "must be a valid UUID",
				Params:  map[string]string{},
			},
		},
		{
			name: "Company not found",
			setup: func() {
				mockCompanyRepository.EXPECT().
					Get(ctx, company.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				id:  company.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CompanyService{
				companyRepository: tt.fields.companyRepository,
				logger:            tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompanyService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyRepository := NewMockcompanyRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var listCompanies []*entity.Company
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		listCompanies = append(listCompanies, mock_models.NewCompany(t))
	}
	filter := mock_models.NewCompanyFilter(t)
	type fields struct {
		companyRepository companyRepository
		logger            log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *entity.CompanyFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*entity.Company
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockCompanyRepository.EXPECT().List(ctx, filter).Return(listCompanies, nil)
				mockCompanyRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listCompanies,
			want1:   count,
			wantErr: nil,
		},
		{
			name:  "invalid",
			setup: func() {},
			fields: fields{
				companyRepository: mockCompanyRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				filter: &entity.CompanyFilter{
					IDs: []entity.UUID{"asd"},
				},
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewInvalidFormError().WithParam("ids", `0: {"code":3,"message":"must be a valid UUID","params":{}}.`),
		},
		{
			name: "list error",
			setup: func() {
				mockCompanyRepository.EXPECT().
					List(ctx, filter).
					Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "count error",
			setup: func() {
				mockCompanyRepository.EXPECT().List(ctx, filter).Return(listCompanies, nil)
				mockCompanyRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CompanyService{
				companyRepository: tt.fields.companyRepository,
				logger:            tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompanyService.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CompanyService.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCompanyService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyRepository := NewMockcompanyRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	ctx := context.Background()
	create := mock_models.NewCompanyCreate(t)
	now := time.Now().UTC()
	type fields struct {
		companyRepository companyRepository
		clock             clock.Clock
		logger            log.Logger
	}
	type args struct {
		ctx    context.Context
		create *entity.CompanyCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *entity.Company
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				mockCompanyRepository.EXPECT().
					Create(
						ctx,
						&entity.Company{
							Name:              create.Name,
							Description:       create.Description,
							AmountOfEmployees: create.AmountOfEmployees,
							Registered:        create.Registered,
							Type:              create.Type,
							UpdatedAt:         now,
							CreatedAt:         now,
						},
					).
					Return(nil)
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &entity.Company{
				ID:                "",
				Name:              create.Name,
				Description:       create.Description,
				AmountOfEmployees: create.AmountOfEmployees,
				Registered:        create.Registered,
				Type:              create.Type,
				UpdatedAt:         now,
				CreatedAt:         now,
			},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				mockCompanyRepository.EXPECT().
					Create(
						ctx,
						&entity.Company{
							ID:                "",
							Name:              create.Name,
							Description:       create.Description,
							AmountOfEmployees: create.AmountOfEmployees,
							Registered:        create.Registered,
							Type:              create.Type,
							UpdatedAt:         now,
							CreatedAt:         now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				create: &entity.CompanyCreate{},
			},
			want: nil,
			wantErr: errs.NewInvalidFormError().WithParams(map[string]string{
				"amount_of_employees": "cannot be blank",
				"name":                "cannot be blank",
				"type":                "cannot be blank",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CompanyService{
				companyRepository: tt.fields.companyRepository,
				clock:             tt.fields.clock,
				logger:            tt.fields.logger,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompanyService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyRepository := NewMockcompanyRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	clockMock := mock_clock.NewMockClock(ctrl)
	update := mock_models.NewCompanyUpdate(t)
	now := company.UpdatedAt
	type fields struct {
		companyRepository companyRepository
		clock             clock.Clock
		logger            log.Logger
	}
	type args struct {
		ctx    context.Context
		update *entity.CompanyUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *entity.Company
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				mockCompanyRepository.EXPECT().
					Get(ctx, update.ID).Return(company, nil)
				mockCompanyRepository.EXPECT().
					Update(ctx, company).Return(nil)
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    company,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				mockCompanyRepository.EXPECT().
					Get(ctx, update.ID).
					Return(company, nil)
				mockCompanyRepository.EXPECT().
					Update(ctx, company).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Company not found",
			setup: func() {
				mockCompanyRepository.EXPECT().Get(ctx, update.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				update: &entity.CompanyUpdate{
					ID: entity.UUID("baduuid"),
				},
			},
			want:    nil,
			wantErr: errs.NewInvalidFormError().WithParam("id", "must be a valid UUID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CompanyService{
				companyRepository: tt.fields.companyRepository,
				clock:             tt.fields.clock,
				logger:            tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompanyService.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyRepository := NewMockcompanyRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	type fields struct {
		companyRepository companyRepository
		logger            log.Logger
	}
	type args struct {
		ctx context.Context
		id  entity.UUID
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockCompanyRepository.EXPECT().
					Delete(ctx, company.ID).
					Return(nil)
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				id:  company.ID,
			},
			wantErr: nil,
		},
		{
			name: "Company not found",
			setup: func() {
				mockCompanyRepository.EXPECT().
					Delete(ctx, company.ID).
					Return(errs.NewEntityNotFound())
			},
			fields: fields{
				companyRepository: mockCompanyRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				id:  company.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CompanyService{
				companyRepository: tt.fields.companyRepository,
				logger:            tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
