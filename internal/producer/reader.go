package producer

import (
	"encoding/json"
	"io"
	"os"

	"github.com/kong/konnect-ingest/internal/models"
	"go.uber.org/zap"
)

// EventReader reads CDC events from a file
type EventReader struct {
	file     *os.File
	decoder  *json.Decoder
	logger   *zap.Logger
}

// NewEventReader creates a new event reader
func NewEventReader(filePath string, logger *zap.Logger) (*EventReader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &EventReader{
		file:    file,
		decoder: json.NewDecoder(file),
		logger:  logger,
	}, nil
}

// ReadEvent reads a single event from the file
func (r *EventReader) ReadEvent() (*models.CDCEvent, error) {
	if !r.decoder.More() {
		return nil, io.EOF
	}

	var event models.CDCEvent
	if err := r.decoder.Decode(&event); err != nil {
		r.logger.Error("Failed to decode event", zap.Error(err))
		return nil, err
	}

	return &event, nil
}

// Close closes the file reader
func (r *EventReader) Close() error {
	return r.file.Close()
}
