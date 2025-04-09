package models

import (
	"encoding/json"
	"fmt"
)

// CDCEvent represents a CDC event from Debezium
type CDCEvent struct {
	Before interface{} `json:"before"`
	After  struct {
		Key   string `json:"key"`
		Value struct {
			Object interface{} `json:"object"`
		} `json:"value"`
	} `json:"after"`
}

// UnmarshalJSON implements json.Unmarshaler
func (e *CDCEvent) UnmarshalJSON(data []byte) error {
	type Alias CDCEvent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Validate required fields
	if aux.After.Key == "" {
		return fmt.Errorf("missing required field: after.key")
	}
	if aux.After.Value.Object == nil {
		return fmt.Errorf("missing required field: after.value.object")
	}

	return nil
}
