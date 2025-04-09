package consumer

import (
	"github.com/Shopify/sarama"
	"github.com/kong/konnect-ingest/internal/models"
)

// EventProcessor defines the contract for processing CDC events
type EventProcessor interface {
	ProcessEvent(event models.CDCEvent) error
}

// MessageHandler defines the contract for handling Kafka messages
type MessageHandler interface {
	sarama.ConsumerGroupHandler
	SetEventProcessor(processor EventProcessor)
}
