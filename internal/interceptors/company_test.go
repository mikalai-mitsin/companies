package interceptors

import (
	"context"
	"errors"
	"github.com/018bf/companies/pkg/utils"
	"reflect"
	"testing"

	"github.com/018bf/companies/internal/domain/errs"
	mock_models "github.com/018bf/companies/internal/domain/models/mock"
	mock_usecases "github.com/018bf/companies/internal/domain/usecases/mock"
	mock_log "github.com/018bf/companies/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"github.com/jaswdr/faker"

	"github.com/018bf/companies/internal/domain/interceptors"
	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/internal/domain/usecases"
	"github.com/018bf/companies/pkg/log"
)

func TestNewCompanyInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	eventUseCase := mock_usecases.NewMockEventUseCase(ctrl)
	companyUseCase := mock_usecases.NewMockCompanyUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase    usecases.AuthUseCase
		companyUseCase usecases.CompanyUseCase
		logger         log.Logger
		eventUseCase   usecases.EventUseCase
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  interceptors.CompanyInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
				eventUseCase:   eventUseCase,
			},
			want: &CompanyInterceptor{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				eventUseCase:   eventUseCase,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewCompanyInterceptor(tt.args.companyUseCase, tt.args.authUseCase, tt.args.eventUseCase, tt.args.logger); !reflect.DeepEqual(
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
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	token := utils.Pointer(mock_models.NewToken(t))
	companyUseCase := mock_usecases.NewMockCompanyUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	type fields struct {
		authUseCase    usecases.AuthUseCase
		companyUseCase usecases.CompanyUseCase
		logger         log.Logger
	}
	type args struct {
		ctx   context.Context
		id    models.UUID
		token *models.Token
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Company
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyDetail).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, company.ID).
					Return(company, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyDetail, company).
					Return(nil)
			},
			fields: fields{
				authUseCase:    authUseCase,
				companyUseCase: companyUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyDetail).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, company.ID).
					Return(company, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyDetail, company).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authUseCase:    authUseCase,
				companyUseCase: companyUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyDetail).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authUseCase:    authUseCase,
				companyUseCase: companyUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyDetail).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, company.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:    authUseCase,
				companyUseCase: companyUseCase,
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
				companyUseCase: tt.fields.companyUseCase,
				authUseCase:    tt.fields.authUseCase,
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
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	token := utils.Pointer(mock_models.NewToken(t))
	companyUseCase := mock_usecases.NewMockCompanyUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	create := mock_models.NewCompanyCreate(t)
	eventUseCase := mock_usecases.NewMockEventUseCase(ctrl)
	type fields struct {
		companyUseCase usecases.CompanyUseCase
		authUseCase    usecases.AuthUseCase
		eventUseCase   usecases.EventUseCase
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.CompanyCreate
		token  *models.Token
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Company
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyCreate, create).
					Return(nil)
				companyUseCase.EXPECT().Create(ctx, create).Return(company, nil)
				eventUseCase.EXPECT().CompanyCreated(ctx, company).Return(nil)
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyCreate, create).
					Return(nil)
				companyUseCase.EXPECT().Create(ctx, create).Return(company, nil)
				eventUseCase.EXPECT().
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
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyCreate, create).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyCreate).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyCreate, create).
					Return(nil)
				companyUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				eventUseCase:   eventUseCase,
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
				companyUseCase: tt.fields.companyUseCase,
				authUseCase:    tt.fields.authUseCase,
				eventUseCase:   tt.fields.eventUseCase,
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
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	token := utils.Pointer(mock_models.NewToken(t))
	companyUseCase := mock_usecases.NewMockCompanyUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	update := mock_models.NewCompanyUpdate(t)
	eventUseCase := mock_usecases.NewMockEventUseCase(ctrl)
	type fields struct {
		companyUseCase usecases.CompanyUseCase
		authUseCase    usecases.AuthUseCase
		eventUseCase   usecases.EventUseCase
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.CompanyUpdate
		token  *models.Token
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Company
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyUpdate).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(company, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyUpdate, company).
					Return(nil)
				companyUseCase.EXPECT().Update(ctx, update).Return(company, nil)
				eventUseCase.EXPECT().CompanyUpdated(ctx, company).Return(nil)
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyUpdate).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(company, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyUpdate, company).
					Return(nil)
				companyUseCase.EXPECT().
					Update(ctx, update).
					Return(company, nil)
				eventUseCase.EXPECT().
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
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyUpdate).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(company, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyUpdate, company).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyUpdate).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyUpdate).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(company, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyUpdate, company).
					Return(nil)
				companyUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyUpdate).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				eventUseCase:   eventUseCase,
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
				companyUseCase: tt.fields.companyUseCase,
				authUseCase:    tt.fields.authUseCase,
				eventUseCase:   tt.fields.eventUseCase,
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
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	token := utils.Pointer(mock_models.NewToken(t))
	companyUseCase := mock_usecases.NewMockCompanyUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	eventUseCase := mock_usecases.NewMockEventUseCase(ctrl)
	type fields struct {
		companyUseCase usecases.CompanyUseCase
		authUseCase    usecases.AuthUseCase
		logger         log.Logger
		eventUseCase   usecases.EventUseCase
	}
	type args struct {
		ctx   context.Context
		id    models.UUID
		token *models.Token
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyDelete).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, company.ID).
					Return(company, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyDelete, company).
					Return(nil)
				companyUseCase.EXPECT().
					Delete(ctx, company.ID).
					Return(nil)
				eventUseCase.EXPECT().
					CompanyDeleted(ctx, company).
					Return(nil)
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyDelete).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, company.ID).
					Return(company, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyDelete, company).
					Return(nil)
				companyUseCase.EXPECT().
					Delete(ctx, company.ID).
					Return(nil)
				eventUseCase.EXPECT().
					CompanyDeleted(ctx, company).
					Return(errs.NewUnexpectedBehaviorError("err 235"))
				logger.EXPECT().Error("can't send 'company deleted' event", log.Context(ctx), log.Error(errs.NewUnexpectedBehaviorError("err 235")))
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyDelete).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, company.ID).
					Return(company, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:    authUseCase,
				companyUseCase: companyUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyDelete).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, company.ID).
					Return(company, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyDelete, company).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authUseCase:    authUseCase,
				companyUseCase: companyUseCase,
				logger:         logger,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyDelete).
					Return(nil)
				companyUseCase.EXPECT().
					Get(ctx, company.ID).
					Return(company, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyDelete, company).
					Return(nil)
				companyUseCase.EXPECT().
					Delete(ctx, company.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:    authUseCase,
				companyUseCase: companyUseCase,
				logger:         logger,
				eventUseCase:   eventUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyDelete).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authUseCase:    authUseCase,
				companyUseCase: companyUseCase,
				logger:         logger,
				eventUseCase:   eventUseCase,
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
				companyUseCase: tt.fields.companyUseCase,
				authUseCase:    tt.fields.authUseCase,
				logger:         tt.fields.logger,
				eventUseCase:   tt.fields.eventUseCase,
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
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	token := utils.Pointer(mock_models.NewToken(t))
	companyUseCase := mock_usecases.NewMockCompanyUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewCompanyFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listCompanies := make([]*models.Company, 0, count)
	for i := uint64(0); i < count; i++ {
		listCompanies = append(listCompanies, mock_models.NewCompany(t))
	}
	type fields struct {
		companyUseCase usecases.CompanyUseCase
		authUseCase    usecases.AuthUseCase
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.CompanyFilter
		token  *models.Token
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Company
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyList, filter).
					Return(nil)
				companyUseCase.EXPECT().
					List(ctx, filter).
					Return(listCompanies, count, nil)
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyList, filter).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyList).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
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
				authUseCase.EXPECT().
					HasPermission(ctx, token, models.PermissionIDCompanyList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, token, models.PermissionIDCompanyList, filter).
					Return(nil)
				companyUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				companyUseCase: companyUseCase,
				authUseCase:    authUseCase,
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
				companyUseCase: tt.fields.companyUseCase,
				authUseCase:    tt.fields.authUseCase,
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
