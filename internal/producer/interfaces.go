package producer

import "github.com/kong/konnect-ingest/internal/models"

// EventProducer defines the contract for producing CDC events
type EventProducer interface {
	ProduceEvent(event models.CDCEvent) error
	Close() error
}
