package data_processing

// DocumentIndexer defines the contract for indexing documents
type DocumentIndexer interface {
	IndexDocument(indexName string, id string, document interface{}) error
}

// EntityExtractor defines the contract for extracting entity information
type EntityExtractor interface {
	ExtractEntityInfo(key string, value interface{}) (entityType string, id string, err error)
}
