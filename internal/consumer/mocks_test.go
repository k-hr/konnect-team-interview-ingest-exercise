package consumer

import (
	"fmt"
)

// MockDocumentIndexer is a mock implementation of DocumentIndexer
type MockDocumentIndexer struct {
	indexedDocs map[string]interface{}
	shouldFail  bool
}

func NewMockDocumentIndexer(shouldFail bool) *MockDocumentIndexer {
	return &MockDocumentIndexer{
		indexedDocs: make(map[string]interface{}),
		shouldFail:  shouldFail,
	}
}

func (m *MockDocumentIndexer) IndexDocument(index, id string, document interface{}) error {
	if m.shouldFail {
		return fmt.Errorf("mock indexing error")
	}
	m.indexedDocs[fmt.Sprintf("%s/%s", index, id)] = document
	return nil
}

// MockEntityExtractor is a mock implementation of EntityExtractor
type MockEntityExtractor struct {
	shouldFail bool
}

func NewMockEntityExtractor(shouldFail bool) *MockEntityExtractor {
	return &MockEntityExtractor{shouldFail: shouldFail}
}

func (m *MockEntityExtractor) ExtractEntityInfo(key string, value interface{}) (string, string, error) {
	if m.shouldFail {
		return "", "", fmt.Errorf("mock extraction error")
	}

	// Simple mock implementation that expects key in format "c/{cluster}/o/{type}/{id}"
	parts := make(map[string]string)
	if obj, ok := value.(map[string]interface{}); ok {
		if id, ok := obj["id"].(string); ok {
			parts["id"] = id
		}
	}

	// Extract type and ID from key
	// Example key: "c/123/o/service/456"
	keyParts := splitKey(key)
	if len(keyParts) >= 5 {
		return keyParts[3], keyParts[4], nil
	}

	return "", "", fmt.Errorf("invalid key format")
}

func splitKey(key string) []string {
	result := make([]string, 0)
	current := ""
	for _, char := range key {
		if char == '/' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
