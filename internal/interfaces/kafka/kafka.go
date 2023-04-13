package kafka

import (
	"fmt"
	"github.com/018bf/companies/internal/configs"
	"github.com/Shopify/sarama"
)

func NewProducer(config *configs.Config) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	server := fmt.Sprintf("%s:%d", config.Kafka.Host, config.Kafka.Port)
	producer, err := sarama.NewSyncProducer([]string{server}, cfg)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
