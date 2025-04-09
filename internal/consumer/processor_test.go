package consumer

import (
	"testing"

	"github.com/kong/konnect-ingest/internal/models"
	"go.uber.org/zap"
)

func TestCDCEventProcessor(t *testing.T) {
	tests := []struct {
		name           string
		event         models.CDCEvent
		indexerFail   bool
		extractorFail bool
		wantErr       bool
	}{
		{
			name: "valid service event",
			event: models.CDCEvent{
				Before: nil,
				After: struct {
					Key   string `json:"key"`
					Value struct {
						Object interface{} `json:"object"`
					} `json:"value"`
				}{
					Key: "c/123/o/service/456",
					Value: struct {
						Object interface{} `json:"object"`
					}{
						Object: map[string]interface{}{
							"name": "test-service",
							"id":   "456",
						},
					},
				},
			},
			indexerFail:   false,
			extractorFail: false,
			wantErr:       false,
		},
		{
			name: "valid node event",
			event: models.CDCEvent{
				Before: nil,
				After: struct {
					Key   string `json:"key"`
					Value struct {
						Object interface{} `json:"object"`
					} `json:"value"`
				}{
					Key: "c/123/o/node/789",
					Value: struct {
						Object interface{} `json:"object"`
					}{
						Object: map[string]interface{}{
							"hostname": "test-node",
							"id":      "789",
						},
					},
				},
			},
			indexerFail:   false,
			extractorFail: false,
			wantErr:       false,
		},
		{
			name: "indexer failure",
			event: models.CDCEvent{
				Before: nil,
				After: struct {
					Key   string `json:"key"`
					Value struct {
						Object interface{} `json:"object"`
					} `json:"value"`
				}{
					Key: "c/123/o/service/456",
					Value: struct {
						Object interface{} `json:"object"`
					}{
						Object: map[string]interface{}{
							"name": "test-service",
							"id":   "456",
						},
					},
				},
			},
			indexerFail:   true,
			extractorFail: false,
			wantErr:       true,
		},
		{
			name: "extractor failure",
			event: models.CDCEvent{
				Before: nil,
				After: struct {
					Key   string `json:"key"`
					Value struct {
						Object interface{} `json:"object"`
					} `json:"value"`
				}{
					Key: "invalid/key",
					Value: struct {
						Object interface{} `json:"object"`
					}{
						Object: map[string]interface{}{
							"name": "test-service",
						},
					},
				},
			},
			indexerFail:   false,
			extractorFail: true,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockIndexer := NewMockDocumentIndexer(tt.indexerFail)
			mockExtractor := NewMockEntityExtractor(tt.extractorFail)

			processor := NewCDCEventProcessor(
				zap.NewNop(),
				mockIndexer,
				mockExtractor,
				"test-index",
			)

			err := processor.ProcessEvent(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("CDCEventProcessor.ProcessEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !tt.indexerFail {
				// Verify the document was indexed correctly
				entityType, id, _ := mockExtractor.ExtractEntityInfo(tt.event.After.Key, tt.event.After.Value.Object)
				indexKey := "test-index-" + entityType + "/" + id
				if _, exists := mockIndexer.indexedDocs[indexKey]; !exists {
					t.Errorf("Document was not indexed with key %s", indexKey)
				}
			}
		})
	}
}
