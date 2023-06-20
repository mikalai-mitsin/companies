package kafka

import (
	"context"
	"encoding/json"
	"github.com/018bf/companies/internal/configs"
	"github.com/018bf/companies/internal/entity"
	"github.com/018bf/companies/pkg/log"
	"github.com/Shopify/sarama"
)

type EventRepository struct {
	producer sarama.SyncProducer
	logger   log.Logger
	topic    string
}

func NewEventRepository(
	producer sarama.SyncProducer,
	config *configs.Config,
	logger log.Logger,
) *EventRepository {
	return &EventRepository{producer: producer, logger: logger, topic: config.Kafka.Topic}
}

func (r *EventRepository) Send(_ context.Context, event *entity.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	message := &sarama.ProducerMessage{
		Topic: r.topic,
		Value: sarama.ByteEncoder(data),
	}
	if _, _, err := r.producer.SendMessage(message); err != nil {
		return err
	}
	return nil
}
