package service

import (
	"context"
	"errors"
	"github.com/018bf/companies/internal/entity"
	mock_models "github.com/018bf/companies/internal/entity/mock"
	"github.com/018bf/companies/internal/errs"
	"github.com/018bf/companies/pkg/log"
	mock_log "github.com/018bf/companies/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestEventService_CompanyCreated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEventRepository := NewMockeventRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	type fields struct {
		eventRepository eventRepository
		logger          log.Logger
	}
	type args struct {
		ctx     context.Context
		company *entity.Company
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
				mockEventRepository.EXPECT().Send(ctx, &entity.Event{
					Operation: entity.EventTypeCreated,
					Company:   company,
				}).Return(nil)
			},
			fields: fields{
				eventRepository: mockEventRepository,
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
				mockEventRepository.EXPECT().
					Send(ctx, &entity.Event{
						Operation: entity.EventTypeCreated,
						Company:   company,
					}).
					Return(errs.NewUnexpectedBehaviorError("err 24"))
			},
			fields: fields{
				eventRepository: mockEventRepository,
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
			u := &EventService{
				eventRepository: tt.fields.eventRepository,
				logger:          tt.fields.logger,
			}
			if err := u.CompanyCreated(tt.args.ctx, tt.args.company); !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyCreated() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventService_CompanyDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEventRepository := NewMockeventRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	type fields struct {
		eventRepository eventRepository
		logger          log.Logger
	}
	type args struct {
		ctx     context.Context
		company *entity.Company
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
				mockEventRepository.EXPECT().Send(ctx, &entity.Event{
					Operation: entity.EventTypeDeleted,
					Company:   company,
				}).Return(nil)
			},
			fields: fields{
				eventRepository: mockEventRepository,
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
				mockEventRepository.EXPECT().
					Send(ctx, &entity.Event{
						Operation: entity.EventTypeDeleted,
						Company:   company,
					}).
					Return(errs.NewUnexpectedBehaviorError("err 24"))
			},
			fields: fields{
				eventRepository: mockEventRepository,
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
			u := &EventService{
				eventRepository: tt.fields.eventRepository,
				logger:          tt.fields.logger,
			}
			if err := u.CompanyDeleted(tt.args.ctx, tt.args.company); !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventService_CompanyUpdated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEventRepository := NewMockeventRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	company := mock_models.NewCompany(t)
	type fields struct {
		eventRepository eventRepository
		logger          log.Logger
	}
	type args struct {
		ctx     context.Context
		company *entity.Company
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
				mockEventRepository.EXPECT().Send(ctx, &entity.Event{
					Operation: entity.EventTypeUpdated,
					Company:   company,
				}).Return(nil)
			},
			fields: fields{
				eventRepository: mockEventRepository,
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
				mockEventRepository.EXPECT().
					Send(ctx, &entity.Event{
						Operation: entity.EventTypeUpdated,
						Company:   company,
					}).
					Return(errs.NewUnexpectedBehaviorError("err 24"))
			},
			fields: fields{
				eventRepository: mockEventRepository,
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
			u := &EventService{
				eventRepository: tt.fields.eventRepository,
				logger:          tt.fields.logger,
			}
			if err := u.CompanyUpdated(tt.args.ctx, tt.args.company); !errors.Is(err, tt.wantErr) {
				t.Errorf("CompanyUpdated() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEventService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEventRepository := NewMockeventRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		eventRepository eventRepository
		logger          log.Logger
	}
	tests := []struct {
		name string
		args args
		want *EventService
	}{
		{
			name: "ok",
			args: args{
				eventRepository: mockEventRepository,
				logger:          logger,
			},
			want: &EventService{
				eventRepository: mockEventRepository,
				logger:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventService(tt.args.eventRepository, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventService() = %v, want %v", got, tt.want)
			}
		})
	}
}
