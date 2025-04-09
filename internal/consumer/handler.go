package consumer

import (
	"github.com/Shopify/sarama"
	"github.com/kong/konnect-ingest/internal/models"
	"go.uber.org/zap"
)

// KafkaConsumerHandler implements MessageHandler for Kafka consumer group
type KafkaConsumerHandler struct {
	logger    *zap.Logger
	processor EventProcessor
}

// NewKafkaConsumerHandler creates a new Kafka consumer handler
func NewKafkaConsumerHandler(logger *zap.Logger) *KafkaConsumerHandler {
	return &KafkaConsumerHandler{
		logger: logger,
	}
}

// SetEventProcessor sets the event processor for the handler
func (h *KafkaConsumerHandler) SetEventProcessor(processor EventProcessor) {
	h.processor = processor
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *KafkaConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *KafkaConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (h *KafkaConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			var event models.CDCEvent
			if err := event.UnmarshalJSON(message.Value); err != nil {
				h.logger.Error("Failed to unmarshal event", zap.Error(err))
				session.MarkMessage(message, "")
				continue
			}

			if err := h.processor.ProcessEvent(event); err != nil {
				h.logger.Error("Failed to process event", zap.Error(err))
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}