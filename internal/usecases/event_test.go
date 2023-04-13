package usecases

import (
	"context"
	"errors"
	"github.com/018bf/companies/internal/domain/errs"
	"github.com/018bf/companies/internal/domain/models"
	mock_models "github.com/018bf/companies/internal/domain/models/mock"
	"github.com/018bf/companies/internal/domain/repositories"
	mock_repositories "github.com/018bf/companies/internal/domain/repositories/mock"
	"github.com/018bf/companies/internal/domain/usecases"
	"github.com/018bf/companies/pkg/log"
	mock_log "github.com/018bf/companies/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestEventUseCase_CompanyCreated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	eventRepository := mock_repositories.NewMockEventRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	type fields struct {
		eventRepository repositories.EventRepository
		logger          log.Logger
	}
	type args struct {
		ctx     context.Context
		company *models.Company
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
				eventRepository.EXPECT().Send(ctx, &models.Event{
					Operation: models.EventTypeCreated,
					Company:   company,
				}).Return(nil)
			},
			fields: fields{
				eventRepository: eventRepository,
				logger:          logger,
			},
			args: args{
				ctx:     ctx,
				company: company,
			},
			wantErr: nil,
		},
		{
			name: "error",
			setup: func() {
				eventRepository.EXPECT().
					Send(ctx, &models.Event{
						Operation: models.EventTypeCreated,
						Company:   company,
					}).
					Return(errs.NewUnexpectedBehaviorError("err 24"))
			},
			fields: fields{
				eventRepository: eventRepository,
				logger:          logger,
			},
			args: args{
				ctx:     ctx,
				company: company,
			},
			wantErr: errs.NewUnexpectedBehaviorError("err 24"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &EventUseCase{
				eventRepository: tt.fields.eventRepository,
				logger:          tt.fields.logger,
			}
			if err := u.CompanyCreated(tt.args.ctx, tt.args.company); !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyCreated() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventUseCase_CompanyDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	eventRepository := mock_repositories.NewMockEventRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	type fields struct {
		eventRepository repositories.EventRepository
		logger          log.Logger
	}
	type args struct {
		ctx     context.Context
		company *models.Company
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
				eventRepository.EXPECT().Send(ctx, &models.Event{
					Operation: models.EventTypeDeleted,
					Company:   company,
				}).Return(nil)
			},
			fields: fields{
				eventRepository: eventRepository,
				logger:          logger,
			},
			args: args{
				ctx:     ctx,
				company: company,
			},
			wantErr: nil,
		},
		{
			name: "error",
			setup: func() {
				eventRepository.EXPECT().
					Send(ctx, &models.Event{
						Operation: models.EventTypeDeleted,
						Company:   company,
					}).
					Return(errs.NewUnexpectedBehaviorError("err 24"))
			},
			fields: fields{
				eventRepository: eventRepository,
				logger:          logger,
			},
			args: args{
				ctx:     ctx,
				company: company,
			},
			wantErr: errs.NewUnexpectedBehaviorError("err 24"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &EventUseCase{
				eventRepository: tt.fields.eventRepository,
				logger:          tt.fields.logger,
			}
			if err := u.CompanyDeleted(tt.args.ctx, tt.args.company); !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventUseCase_CompanyUpdated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	eventRepository := mock_repositories.NewMockEventRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	type fields struct {
		eventRepository repositories.EventRepository
		logger          log.Logger
	}
	type args struct {
		ctx     context.Context
		company *models.Company
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
				eventRepository.EXPECT().Send(ctx, &models.Event{
					Operation: models.EventTypeUpdated,
					Company:   company,
				}).Return(nil)
			},
			fields: fields{
				eventRepository: eventRepository,
				logger:          logger,
			},
			args: args{
				ctx:     ctx,
				company: company,
			},
			wantErr: nil,
		},
		{
			name: "error",
			setup: func() {
				eventRepository.EXPECT().
					Send(ctx, &models.Event{
						Operation: models.EventTypeUpdated,
						Company:   company,
					}).
					Return(errs.NewUnexpectedBehaviorError("err 24"))
			},
			fields: fields{
				eventRepository: eventRepository,
				logger:          logger,
			},
			args: args{
				ctx:     ctx,
				company: company,
			},
			wantErr: errs.NewUnexpectedBehaviorError("err 24"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &EventUseCase{
				eventRepository: tt.fields.eventRepository,
				logger:          tt.fields.logger,
			}
			if err := u.CompanyUpdated(tt.args.ctx, tt.args.company); !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyUpdated() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEventUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	eventRepository := mock_repositories.NewMockEventRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		eventRepository repositories.EventRepository
		logger          log.Logger
	}
	tests := []struct {
		name string
		args args
		want usecases.EventUseCase
	}{
		{
			name: "ok",
			args: args{
				eventRepository: eventRepository,
				logger:          logger,
			},
			want: &EventUseCase{
				eventRepository: eventRepository,
				logger:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventUseCase(tt.args.eventRepository, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
