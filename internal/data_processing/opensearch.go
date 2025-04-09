package data_processing

import (
	"encoding/json"
	"strings"

	"github.com/opensearch-project/opensearch-go/v2"
	"go.uber.org/zap"
)

// OpenSearchIndexer implements DocumentIndexer for OpenSearch
type OpenSearchIndexer struct {
	client *opensearch.Client
	logger *zap.Logger
}

// NewOpenSearchIndexer creates a new OpenSearch indexer
func NewOpenSearchIndexer(client *opensearch.Client, logger *zap.Logger) *OpenSearchIndexer {
	return &OpenSearchIndexer{
		client: client,
		logger: logger,
	}
}

// IndexDocument indexes a document in OpenSearch
func (i *OpenSearchIndexer) IndexDocument(indexName string, id string, document interface{}) error {
	objectBytes, err := json.Marshal(document)
	if err != nil {
		i.logger.Error("Failed to marshal object", zap.Error(err))
		return err
	}

	_, err = i.client.Index(
		indexName,
		strings.NewReader(string(objectBytes)),
		i.client.Index.WithDocumentID(id),
	)
	if err != nil {
		i.logger.Error("Failed to index document", zap.Error(err))
		return err
	}

	return nil
}
