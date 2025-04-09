package consumer

import (
	"github.com/kong/konnect-ingest/internal/data_processing"
	"github.com/kong/konnect-ingest/internal/models"
	"go.uber.org/zap"
)

// CDCEventProcessor implements EventProcessor for CDC events
type CDCEventProcessor struct {
	logger          *zap.Logger
	indexer         data_processing.DocumentIndexer
	entityExtractor data_processing.EntityExtractor
	indexPrefix     string
}

// NewCDCEventProcessor creates a new CDC event processor
func NewCDCEventProcessor(
	logger *zap.Logger,
	indexer data_processing.DocumentIndexer,
	entityExtractor data_processing.EntityExtractor,
	indexPrefix string,
) *CDCEventProcessor {
	return &CDCEventProcessor{
		logger:          logger,
		indexer:         indexer,
		entityExtractor: entityExtractor,
		indexPrefix:     indexPrefix,
	}
}

// ProcessEvent processes a single CDC event
func (p *CDCEventProcessor) ProcessEvent(event models.CDCEvent) error {
	// Extract entity type and ID
	entityType, id, err := p.entityExtractor.ExtractEntityInfo(event.After.Key, event.After.Value.Object)
	if err != nil {
		return err
	}

	// Create index name based on entity type
	indexName := p.indexPrefix + "-" + entityType
	p.logger.Info("Processing event",
		zap.String("entityType", entityType),
		zap.String("indexName", indexName),
		zap.String("key", event.After.Key),
	)

	// Index the document
	if err := p.indexer.IndexDocument(indexName, id, event.After.Value.Object); err != nil {
		return err
	}

	p.logger.Info("Successfully indexed document",
		zap.String("indexName", indexName),
		zap.String("id", id),
	)

	return nil
}