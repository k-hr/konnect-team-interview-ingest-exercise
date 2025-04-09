package models

import (
	"encoding/json"
	"testing"
)

func TestCDCEventUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    CDCEvent
		wantErr bool
	}{
		{
			name: "valid service event",
			json: `{
				"before": null,
				"after": {
					"key": "c/123/o/service/456",
					"value": {
						"object": {"name": "test-service", "id": "456"}
					}
				}
			}`,
			want: CDCEvent{
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
			wantErr: false,
		},
		{
			name: "valid node event",
			json: `{
				"before": {"id": "old-id"},
				"after": {
					"key": "c/123/o/node/789",
					"value": {
						"object": {"hostname": "test-node", "id": "789"}
					}
				}
			}`,
			want: CDCEvent{
				Before: map[string]interface{}{"id": "old-id"},
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
			wantErr: false,
		},
		{
			name:    "invalid json",
			json:    `{invalid`,
			wantErr: true,
		},
		{
			name: "missing required fields",
			json: `{
				"before": null
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got CDCEvent
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("CDCEvent.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				gotJSON, err := json.Marshal(got)
				if err != nil {
					t.Errorf("Failed to marshal result: %v", err)
					return
				}

				wantJSON, err := json.Marshal(tt.want)
				if err != nil {
					t.Errorf("Failed to marshal expected: %v", err)
					return
				}

				if string(gotJSON) != string(wantJSON) {
					t.Errorf("CDCEvent.UnmarshalJSON() = %v, want %v", string(gotJSON), string(wantJSON))
				}
			}
		})
	}
}
