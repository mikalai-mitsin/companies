package kafka

import (
	"context"
	"errors"
	"github.com/018bf/companies/internal/configs"
	"github.com/018bf/companies/internal/entity"
	mock_models "github.com/018bf/companies/internal/entity/mock"
	"github.com/018bf/companies/internal/errs"
	mock_sarama "github.com/018bf/companies/internal/event/repositories/kafka/mock"
	"github.com/018bf/companies/pkg/log"
	mock_log "github.com/018bf/companies/pkg/log/mock"
	"github.com/Shopify/sarama"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

//go:generate mockgen -build_flags=-mod=mod -destination mock/sarma.go github.com/Shopify/sarama SyncProducer

func TestEventRepository_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	syncProducer := mock_sarama.NewMockSyncProducer(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	config := configs.NewMockConfig(t)
	ctx := context.Background()
	event := mock_models.NewEvent(t)
	type fields struct {
		producer sarama.SyncProducer
		logger   log.Logger
		topic    string
	}
	type args struct {
		in0   context.Context
		event *entity.Event
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
				syncProducer.EXPECT().SendMessage(gomock.Any()).Return(int32(0), int64(0), nil)
			},
			fields: fields{
				producer: syncProducer,
				logger:   logger,
				topic:    config.Kafka.Topic,
			},
			args: args{
				in0:   ctx,
				event: event,
			},
			wantErr: nil,
		},
		{
			name: "send error",
			setup: func() {
				syncProducer.EXPECT().SendMessage(gomock.Any()).Return(int32(0), int64(0), errs.NewUnexpectedBehaviorError("err 1234"))
			},
			fields: fields{
				producer: syncProducer,
				logger:   logger,
				topic:    config.Kafka.Topic,
			},
			args: args{
				in0:   ctx,
				event: event,
			},
			wantErr: errs.NewUnexpectedBehaviorError("err 1234"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &EventRepository{
				producer: tt.fields.producer,
				logger:   tt.fields.logger,
				topic:    tt.fields.topic,
			}
			if err := r.Send(tt.args.in0, tt.args.event); !errors.Is(err, tt.wantErr) {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEventRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	syncProducer := mock_sarama.NewMockSyncProducer(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	config := configs.NewMockConfig(t)
	type args struct {
		producer sarama.SyncProducer
		config   *configs.Config
		logger   log.Logger
	}
	tests := []struct {
		name string
		args args
		want *EventRepository
	}{
		{
			name: "ok",
			args: args{
				producer: syncProducer,
				config:   config,
				logger:   logger,
			},
			want: &EventRepository{
				producer: syncProducer,
				topic:    config.Kafka.Topic,
				logger:   logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventRepository(tt.args.producer, tt.args.config, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
