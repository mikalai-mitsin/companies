package interceptor

import (
	"context"
	"errors"
	"github.com/018bf/companies/pkg/utils"
	"reflect"
	"testing"

	mock_models "github.com/018bf/companies/internal/entity/mock"
	"github.com/018bf/companies/internal/errs"
	mock_log "github.com/018bf/companies/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"github.com/jaswdr/faker"

	"github.com/018bf/companies/internal/entity"
	"github.com/018bf/companies/pkg/log"
)

func TestNewCompanyInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := NewMockauthService(ctrl)
	mockEventService := NewMockeventService(ctrl)
	mockCompanyService := NewMockcompanyService(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authService    authService
		companyService companyService
		logger         log.Logger
		eventService   eventService
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *CompanyInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				logger:         logger,
				eventService:   mockEventService,
			},
			want: &CompanyInterceptor{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				eventService:   mockEventService,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewCompanyInterceptor(tt.args.companyService, tt.args.authService, tt.args.eventService, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewCompanyInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := NewMockauthService(ctrl)
	token := utils.Pointer(mock_models.NewToken(t))
	mockCompanyService := NewMockcompanyService(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	type fields struct {
		authService    authService
		companyService companyService
		logger         log.Logger
	}
	type args struct {
		ctx   context.Context
		id    entity.UUID
		token *entity.Token
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
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyDetail).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, company.ID).
					Return(company, nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyDetail, company).
					Return(nil)
			},
			fields: fields{
				authService:    mockAuthService,
				companyService: mockCompanyService,
				logger:         logger,
			},
			args: args{
				ctx:   ctx,
				id:    company.ID,
				token: token,
			},
			want:    company,
			wantErr: nil,
		},
		{
			name: "object permission error",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyDetail).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, company.ID).
					Return(company, nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyDetail, company).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authService:    mockAuthService,
				companyService: mockCompanyService,
				logger:         logger,
			},
			args: args{
				ctx:   ctx,
				id:    company.ID,
				token: token,
			},
			want:    nil,
			wantErr: errs.NewPermissionDenied(),
		},
		{
			name: "permission denied",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyDetail).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authService:    mockAuthService,
				companyService: mockCompanyService,
				logger:         logger,
			},
			args: args{
				ctx:   ctx,
				id:    company.ID,
				token: token,
			},
			want:    nil,
			wantErr: errs.NewPermissionDenied(),
		},
		{
			name: "Company not found",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyDetail).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, company.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authService:    mockAuthService,
				companyService: mockCompanyService,
				logger:         logger,
			},
			args: args{
				ctx:   ctx,
				id:    company.ID,
				token: token,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &CompanyInterceptor{
				companyService: tt.fields.companyService,
				authService:    tt.fields.authService,
				logger:         tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id, tt.args.token)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompanyInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := NewMockauthService(ctrl)
	token := utils.Pointer(mock_models.NewToken(t))
	mockCompanyService := NewMockcompanyService(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	create := mock_models.NewCompanyCreate(t)
	mockEventService := NewMockeventService(ctrl)
	type fields struct {
		companyService companyService
		authService    authService
		eventService   eventService
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		create *entity.CompanyCreate
		token  *entity.Token
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
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyCreate).
					Return(nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyCreate, create).
					Return(nil)
				mockCompanyService.EXPECT().Create(ctx, create).Return(company, nil)
				mockEventService.EXPECT().CompanyCreated(ctx, company).Return(nil)
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				eventService:   mockEventService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
				token:  token,
			},
			want:    company,
			wantErr: nil,
		},
		{
			name: "send event error",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyCreate).
					Return(nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyCreate, create).
					Return(nil)
				mockCompanyService.EXPECT().Create(ctx, create).Return(company, nil)
				mockEventService.EXPECT().
					CompanyCreated(ctx, company).
					Return(errs.NewUnexpectedBehaviorError("err 235"))
				logger.EXPECT().
					Error(
						"can't send 'company created' event",
						log.Context(ctx),
						log.Error(errs.NewUnexpectedBehaviorError("err 235")),
					)
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				eventService:   mockEventService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
				token:  token,
			},
			want:    company,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyCreate).
					Return(nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyCreate, create).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				eventService:   mockEventService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
				token:  token,
			},
			want:    nil,
			wantErr: errs.NewPermissionDenied(),
		},
		{
			name: "permission denied",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyCreate).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				eventService:   mockEventService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
				token:  token,
			},
			want:    nil,
			wantErr: errs.NewPermissionDenied(),
		},
		{
			name: "create error",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyCreate).
					Return(nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyCreate, create).
					Return(nil)
				mockCompanyService.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				eventService:   mockEventService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
				token:  token,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &CompanyInterceptor{
				companyService: tt.fields.companyService,
				authService:    tt.fields.authService,
				eventService:   tt.fields.eventService,
				logger:         tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create, tt.args.token)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompanyInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := NewMockauthService(ctrl)
	token := utils.Pointer(mock_models.NewToken(t))
	mockCompanyService := NewMockcompanyService(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	update := mock_models.NewCompanyUpdate(t)
	mockEventService := NewMockeventService(ctrl)
	type fields struct {
		companyService companyService
		authService    authService
		eventService   eventService
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		update *entity.CompanyUpdate
		token  *entity.Token
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
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyUpdate).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, update.ID).
					Return(company, nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyUpdate, company).
					Return(nil)
				mockCompanyService.EXPECT().Update(ctx, update).Return(company, nil)
				mockEventService.EXPECT().CompanyUpdated(ctx, company).Return(nil)
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				eventService:   mockEventService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
				token:  token,
			},
			want:    company,
			wantErr: nil,
		},
		{
			name: "send event error",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyUpdate).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, update.ID).
					Return(company, nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyUpdate, company).
					Return(nil)
				mockCompanyService.EXPECT().
					Update(ctx, update).
					Return(company, nil)
				mockEventService.EXPECT().
					CompanyUpdated(ctx, company).
					Return(errs.NewUnexpectedBehaviorError("err 235"))
				logger.EXPECT().
					Error(
						"can't send 'company updated' event",
						log.Context(ctx),
						log.Error(errs.NewUnexpectedBehaviorError("err 235")),
					)
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				eventService:   mockEventService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
				token:  token,
			},
			want:    company,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyUpdate).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, update.ID).
					Return(company, nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyUpdate, company).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				eventService:   mockEventService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
				token:  token,
			},
			want:    nil,
			wantErr: errs.NewPermissionDenied(),
		},
		{
			name: "not found",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyUpdate).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				eventService:   mockEventService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
				token:  token,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "update error",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyUpdate).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, update.ID).
					Return(company, nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyUpdate, company).
					Return(nil)
				mockCompanyService.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				eventService:   mockEventService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
				token:  token,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyUpdate).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				eventService:   mockEventService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
				token:  token,
			},
			wantErr: errs.NewPermissionDenied(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &CompanyInterceptor{
				companyService: tt.fields.companyService,
				authService:    tt.fields.authService,
				eventService:   tt.fields.eventService,
				logger:         tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update, tt.args.token)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompanyInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompanyInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := NewMockauthService(ctrl)
	token := utils.Pointer(mock_models.NewToken(t))
	mockCompanyService := NewMockcompanyService(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	mockEventService := NewMockeventService(ctrl)
	type fields struct {
		companyService companyService
		authService    authService
		logger         log.Logger
		eventService   eventService
	}
	type args struct {
		ctx   context.Context
		id    entity.UUID
		token *entity.Token
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
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyDelete).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, company.ID).
					Return(company, nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyDelete, company).
					Return(nil)
				mockCompanyService.EXPECT().
					Delete(ctx, company.ID).
					Return(nil)
				mockEventService.EXPECT().
					CompanyDeleted(ctx, company).
					Return(nil)
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				logger:         logger,
				eventService:   mockEventService,
			},
			args: args{
				ctx:   ctx,
				id:    company.ID,
				token: token,
			},
			wantErr: nil,
		},
		{
			name: "send event error",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyDelete).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, company.ID).
					Return(company, nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyDelete, company).
					Return(nil)
				mockCompanyService.EXPECT().
					Delete(ctx, company.ID).
					Return(nil)
				mockEventService.EXPECT().
					CompanyDeleted(ctx, company).
					Return(errs.NewUnexpectedBehaviorError("err 235"))
				logger.EXPECT().Error("can't send 'company deleted' event", log.Context(ctx), log.Error(errs.NewUnexpectedBehaviorError("err 235")))
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				logger:         logger,
				eventService:   mockEventService,
			},
			args: args{
				ctx:   ctx,
				id:    company.ID,
				token: token,
			},
			wantErr: nil,
		},
		{
			name: "Company not found",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyDelete).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, company.ID).
					Return(company, errs.NewEntityNotFound())
			},
			fields: fields{
				authService:    mockAuthService,
				companyService: mockCompanyService,
				logger:         logger,
			},
			args: args{
				ctx:   ctx,
				id:    company.ID,
				token: token,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyDelete).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, company.ID).
					Return(company, nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyDelete, company).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authService:    mockAuthService,
				companyService: mockCompanyService,
				logger:         logger,
				eventService:   mockEventService,
			},
			args: args{
				ctx:   ctx,
				id:    company.ID,
				token: token,
			},
			wantErr: errs.NewPermissionDenied(),
		},
		{
			name: "delete error",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyDelete).
					Return(nil)
				mockCompanyService.EXPECT().
					Get(ctx, company.ID).
					Return(company, nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyDelete, company).
					Return(nil)
				mockCompanyService.EXPECT().
					Delete(ctx, company.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authService:    mockAuthService,
				companyService: mockCompanyService,
				logger:         logger,
				eventService:   mockEventService,
			},
			args: args{
				ctx:   ctx,
				id:    company.ID,
				token: token,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyDelete).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authService:    mockAuthService,
				companyService: mockCompanyService,
				logger:         logger,
				eventService:   mockEventService,
			},
			args: args{
				ctx:   ctx,
				id:    company.ID,
				token: token,
			},
			wantErr: errs.NewPermissionDenied(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &CompanyInterceptor{
				companyService: tt.fields.companyService,
				authService:    tt.fields.authService,
				logger:         tt.fields.logger,
				eventService:   tt.fields.eventService,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id, tt.args.token); !errors.Is(
				err,
				tt.wantErr,
			) {
				t.Errorf("CompanyInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCompanyInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := NewMockauthService(ctrl)
	token := utils.Pointer(mock_models.NewToken(t))
	mockCompanyService := NewMockcompanyService(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewCompanyFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listCompanies := make([]*entity.Company, 0, count)
	for i := uint64(0); i < count; i++ {
		listCompanies = append(listCompanies, mock_models.NewCompany(t))
	}
	type fields struct {
		companyService companyService
		authService    authService
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *entity.CompanyFilter
		token  *entity.Token
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
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyList).
					Return(nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyList, filter).
					Return(nil)
				mockCompanyService.EXPECT().
					List(ctx, filter).
					Return(listCompanies, count, nil)
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
				token:  token,
			},
			want:    listCompanies,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyList).
					Return(nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyList, filter).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
				token:  token,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDenied(),
		},
		{
			name: "permission error",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyList).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
				token:  token,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDenied(),
		},
		{
			name: "list error",
			setup: func() {
				mockAuthService.EXPECT().
					HasPermission(ctx, token, entity.PermissionIDCompanyList).
					Return(nil)
				mockAuthService.EXPECT().
					HasObjectPermission(ctx, token, entity.PermissionIDCompanyList, filter).
					Return(nil)
				mockCompanyService.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				companyService: mockCompanyService,
				authService:    mockAuthService,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
				token:  token,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("l e"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &CompanyInterceptor{
				companyService: tt.fields.companyService,
				authService:    tt.fields.authService,
				logger:         tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter, tt.args.token)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompanyInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CompanyInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
