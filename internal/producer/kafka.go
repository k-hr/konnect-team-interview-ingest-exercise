package producer

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/kong/konnect-ingest/internal/models"
	"go.uber.org/zap"
)

// KafkaEventProducer implements EventProducer for Kafka
type KafkaEventProducer struct {
	producer sarama.SyncProducer
	topic    string
	logger   *zap.Logger
}

// NewKafkaEventProducer creates a new Kafka event producer
func NewKafkaEventProducer(brokers []string, topic string, logger *zap.Logger) (*KafkaEventProducer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewSyncProducer(brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaEventProducer{
		producer: producer,
		topic:    topic,
		logger:   logger,
	}, nil
}

// ProduceEvent produces a single CDC event to Kafka
func (p *KafkaEventProducer) ProduceEvent(event models.CDCEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Failed to marshal event", zap.Error(err))
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(event.After.Key),
		Value: sarama.StringEncoder(eventBytes),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.Error("Failed to send message", zap.Error(err))
		return err
	}

	p.logger.Info("Message sent",
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)

	return nil
}

// Close closes the Kafka producer
func (p *KafkaEventProducer) Close() error {
	return p.producer.Close()
}
