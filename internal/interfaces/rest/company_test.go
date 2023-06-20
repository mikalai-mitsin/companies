package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/018bf/companies/pkg/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/018bf/companies/internal/entity"
	mock_models "github.com/018bf/companies/internal/entity/mock"
	"github.com/018bf/companies/internal/errs"
	"github.com/018bf/companies/pkg/log"
	mock_log "github.com/018bf/companies/pkg/log/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestCompanyHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyInterceptor := NewMockcompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	create := mock_models.NewCompanyCreate(t)
	createjson, _ := json.Marshal(create)
	company := mock_models.NewCompany(t)
	companyjson, _ := json.Marshal(company)
	type fields struct {
		companyInterceptor companyInterceptor
		logger             log.Logger
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name       string
		setup      func()
		fields     fields
		args       args
		wantBody   *bytes.Buffer
		wantStatus int
	}{
		{
			name: "ok",
			setup: func() {
				mockCompanyInterceptor.EXPECT().
					Create(gomock.Any(), create, utils.Pointer(entity.Token("good token"))).
					Return(company, nil)
			},
			fields: fields{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(createjson)),
				}).WithContext(context.WithValue(context.Background(), TokenContextKey, utils.Pointer(entity.Token("good token")))),
			},
			wantBody:   bytes.NewBuffer(companyjson),
			wantStatus: http.StatusCreated,
		},
		{
			name: "permission denied",
			setup: func() {
				mockCompanyInterceptor.EXPECT().
					Create(gomock.Any(), create, utils.Pointer(entity.Token("good token"))).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(createjson)),
				}).WithContext(context.WithValue(context.Background(), TokenContextKey, utils.Pointer(entity.Token("good token")))),
			},
			wantBody:   bytes.NewBufferString(errs.NewPermissionDenied().Error()),
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			h := &CompanyHandler{
				companyInterceptor: tt.fields.companyInterceptor,
				logger:             tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			h.Create(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("Create() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("Create() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestCompanyHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyInterceptor := NewMockcompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	company := mock_models.NewCompany(t)
	type fields struct {
		companyInterceptor companyInterceptor
		logger             log.Logger
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name       string
		setup      func()
		fields     fields
		args       args
		wantStatus int
		wantBody   *bytes.Buffer
	}{
		{
			name: "ok",
			setup: func() {
				mockCompanyInterceptor.EXPECT().
					Delete(gomock.Any(), company.ID, utils.Pointer(entity.Token("good token"))).
					Return(nil)
			},
			fields: fields{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), TokenContextKey, utils.Pointer(entity.Token("good token"))),
				),
			},
			wantBody:   &bytes.Buffer{},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "permission denied",
			setup: func() {
				mockCompanyInterceptor.EXPECT().
					Delete(gomock.Any(), company.ID, utils.Pointer(entity.Token("good token"))).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), TokenContextKey, utils.Pointer(entity.Token("good token"))),
				),
			},
			wantBody:   bytes.NewBufferString(errs.NewPermissionDenied().Error()),
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			h := &CompanyHandler{
				companyInterceptor: tt.fields.companyInterceptor,
				logger:             tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(company.ID))
			h.Delete(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("Delete() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("Delete() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestCompanyHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyInterceptor := NewMockcompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	company := mock_models.NewCompany(t)
	companyjson, _ := json.Marshal(company)
	type fields struct {
		companyInterceptor companyInterceptor
		logger             log.Logger
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name       string
		setup      func()
		fields     fields
		args       args
		wantStatus int
		wantBody   *bytes.Buffer
	}{
		{
			name: "ok",
			setup: func() {
				mockCompanyInterceptor.EXPECT().
					Get(gomock.Any(), company.ID, utils.Pointer(entity.Token("good token"))).
					Return(company, nil)
			},
			fields: fields{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), TokenContextKey, utils.Pointer(entity.Token("good token"))),
				),
			},
			wantBody:   bytes.NewBuffer(companyjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				mockCompanyInterceptor.EXPECT().
					Get(gomock.Any(), company.ID, utils.Pointer(entity.Token("good token"))).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), TokenContextKey, utils.Pointer(entity.Token("good token"))),
				),
			},
			wantBody:   bytes.NewBufferString(errs.NewPermissionDenied().Error()),
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			h := &CompanyHandler{
				companyInterceptor: tt.fields.companyInterceptor,
				logger:             tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(company.ID))
			h.Get(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("Get() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("Get() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestCompanyHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyInterceptor := NewMockcompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	filter := &entity.CompanyFilter{}
	listCompanies := []*entity.Company{mock_models.NewCompany(t)}
	listCompaniesjson, _ := json.Marshal(listCompanies)
	type fields struct {
		companyInterceptor companyInterceptor
		logger             log.Logger
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name       string
		setup      func()
		fields     fields
		args       args
		wantStatus int
		wantBody   *bytes.Buffer
	}{
		{
			name: "ok",
			setup: func() {
				mockCompanyInterceptor.EXPECT().
					List(gomock.Any(), filter, utils.Pointer(entity.Token("good token"))).
					Return(listCompanies, uint64(len(listCompanies)), nil)
			},
			fields: fields{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), TokenContextKey, utils.Pointer(entity.Token("good token"))),
				),
			},
			wantBody:   bytes.NewBuffer(listCompaniesjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				mockCompanyInterceptor.EXPECT().
					List(gomock.Any(), filter, utils.Pointer(entity.Token("good token"))).
					Return(nil, uint64(0), errs.NewPermissionDenied())
			},
			fields: fields{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), TokenContextKey, utils.Pointer(entity.Token("good token"))),
				),
			},
			wantBody:   bytes.NewBufferString(errs.NewPermissionDenied().Error()),
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			h := &CompanyHandler{
				companyInterceptor: tt.fields.companyInterceptor,
				logger:             tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			h.List(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("List() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("List() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestCompanyHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyInterceptor := NewMockcompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type fields struct {
		companyInterceptor companyInterceptor
		logger             log.Logger
	}
	type args struct {
		router *gin.RouterGroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ok",
			fields: fields{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
			args: args{
				router: gin.Default().Group("/"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CompanyHandler{
				companyInterceptor: tt.fields.companyInterceptor,
				logger:             tt.fields.logger,
			}
			h.Register(tt.args.router)
		})
	}
}

func TestCompanyHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyInterceptor := NewMockcompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	update := mock_models.NewCompanyUpdate(t)
	updatejson, _ := json.Marshal(update)
	company := mock_models.NewCompany(t)
	companyjson, _ := json.Marshal(company)
	type fields struct {
		companyInterceptor companyInterceptor
		logger             log.Logger
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name       string
		setup      func()
		fields     fields
		args       args
		wantStatus int
		wantBody   *bytes.Buffer
	}{

		{
			name: "ok",
			setup: func() {
				mockCompanyInterceptor.EXPECT().
					Update(gomock.Any(), update, utils.Pointer(entity.Token("good token"))).
					Return(company, nil)
			},
			fields: fields{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(updatejson)),
				}).WithContext(context.WithValue(context.Background(), TokenContextKey, utils.Pointer(entity.Token("good token")))),
			},
			wantBody:   bytes.NewBuffer(companyjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				mockCompanyInterceptor.EXPECT().
					Update(gomock.Any(), update, utils.Pointer(entity.Token("good token"))).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(updatejson)),
				}).WithContext(context.WithValue(context.Background(), TokenContextKey, utils.Pointer(entity.Token("good token")))),
			},
			wantBody:   bytes.NewBufferString(errs.NewPermissionDenied().Error()),
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			h := &CompanyHandler{
				companyInterceptor: tt.fields.companyInterceptor,
				logger:             tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(update.ID))
			h.Update(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("Update() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("Update() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestNewCompanyHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCompanyInterceptor := NewMockcompanyInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		companyInterceptor companyInterceptor
		logger             log.Logger
	}
	tests := []struct {
		name string
		args args
		want *CompanyHandler
	}{
		{
			name: "ok",
			args: args{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
			want: &CompanyHandler{
				companyInterceptor: mockCompanyInterceptor,
				logger:             logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCompanyHandler(tt.args.companyInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewCompanyHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
