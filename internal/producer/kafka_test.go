package producer

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/kong/konnect-ingest/internal/models"
	"go.uber.org/zap"
)

// mockSyncProducer implements sarama.SyncProducer for testing
type mockSyncProducer struct {
	sendMessageError error
	messages         []*sarama.ProducerMessage
}

func (m *mockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	if m.sendMessageError != nil {
		return 0, 0, m.sendMessageError
	}
	m.messages = append(m.messages, msg)
	return 0, int64(len(m.messages)), nil
}

func (m *mockSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	for _, msg := range msgs {
		_, _, err := m.SendMessage(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mockSyncProducer) Close() error {
	return nil
}

func (m *mockSyncProducer) TxnStatus() sarama.ProducerTxnStatusFlag {
	return 0
}

func (m *mockSyncProducer) IsTransactional() bool {
	return false
}

func (m *mockSyncProducer) BeginTxn() error {
	return nil
}

func (m *mockSyncProducer) CommitTxn() error {
	return nil
}

func (m *mockSyncProducer) AbortTxn() error {
	return nil
}

func (m *mockSyncProducer) AddMessageToTxn(msg *sarama.ConsumerMessage, groupID string, topic *string) error {
	return nil
}

func (m *mockSyncProducer) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	return nil
}

func TestKafkaEventProducer(t *testing.T) {
	tests := []struct {
		name          string
		event         models.CDCEvent
		producerError error
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
			wantErr: false,
		},
		{
			name: "producer error",
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
			producerError: sarama.ErrBrokerNotAvailable,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock producer for each test
			mockProducer := &mockSyncProducer{sendMessageError: tt.producerError}
			producer := &KafkaEventProducer{
				producer: mockProducer,
				topic:    "test-topic",
				logger:   zap.NewNop(),
			}

			// Test producing event
			err := producer.ProduceEvent(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("KafkaEventProducer.ProduceEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify message was sent correctly if no error expected
			if !tt.wantErr {
				if len(mockProducer.messages) != 1 {
					t.Errorf("Expected 1 message to be sent, got %d", len(mockProducer.messages))
					return
				}

				msg := mockProducer.messages[0]
				if msg.Topic != "test-topic" {
					t.Errorf("Expected topic 'test-topic', got %s", msg.Topic)
				}
				if string(msg.Key.(sarama.StringEncoder)) != tt.event.After.Key {
					t.Errorf("Expected key %s, got %s", tt.event.After.Key, string(msg.Key.(sarama.StringEncoder)))
				}
			}
		})
	}
}
