package data_processing

import (
	"errors"
	"strings"

	"go.uber.org/zap"
)

var (
	ErrInvalidKey    = errors.New("invalid key format")
	ErrInvalidObject = errors.New("invalid object format")
	ErrMissingID     = errors.New("missing id field")
)

// CDCEntityExtractor implements EntityExtractor for CDC events
type CDCEntityExtractor struct {
	logger *zap.Logger
}

// NewCDCEntityExtractor creates a new CDC entity extractor
func NewCDCEntityExtractor(logger *zap.Logger) *CDCEntityExtractor {
	return &CDCEntityExtractor{
		logger: logger,
	}
}

// ExtractEntityInfo extracts entity type and ID from CDC event data
func (e *CDCEntityExtractor) ExtractEntityInfo(key string, value interface{}) (string, string, error) {
	parts := strings.Split(key, "/")
	if len(parts) < 4 {
		e.logger.Error("Key has insufficient parts", zap.String("key", key))
		return "", "", ErrInvalidKey
	}

	entityType := parts[len(parts)-2] // e.g., "service", "node", "upstream"

	obj, ok := value.(map[string]interface{})
	if !ok {
		e.logger.Error("Failed to cast object to map", zap.String("key", key))
		return "", "", ErrInvalidObject
	}

	id, ok := obj["id"].(string)
	if !ok {
		e.logger.Error("Failed to get ID from object", zap.String("key", key))
		return "", "", ErrMissingID
	}

	return entityType, id, nil
}