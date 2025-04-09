package models

import "encoding/json"

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
	return json.Unmarshal(data, aux)
}
