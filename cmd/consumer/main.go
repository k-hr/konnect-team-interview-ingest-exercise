package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/kong/konnect-ingest/internal/config"
	"github.com/kong/konnect-ingest/internal/models"
	"github.com/opensearch-project/opensearch-go/v2"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Create OpenSearch client
	osClient, err := opensearch.NewClient(opensearch.Config{
		Addresses: cfg.OpenSearch.Hosts,
	})
	if err != nil {
		logger.Fatal("Failed to create OpenSearch client", zap.Error(err))
	}

	// Create Kafka consumer
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.GroupID, kafkaConfig)
	if err != nil {
		logger.Fatal("Failed to create consumer group", zap.Error(err))
	}
	defer consumer.Close()

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		logger.Info("Received shutdown signal, stopping...")
		cancel()
	}()

	// Start consuming
	for {
		handler := &ConsumerGroupHandler{
			logger:      logger,
			osClient:    osClient,
			indexPrefix: cfg.OpenSearch.IndexPrefix,
		}

		err := consumer.Consume(ctx, []string{cfg.Kafka.Topic}, handler)
		if err != nil {
			if err == context.Canceled {
				return
			}
			logger.Error("Error from consumer", zap.Error(err))
		}
	}
}

// ConsumerGroupHandler handles the consumer group session
type ConsumerGroupHandler struct {
	logger      *zap.Logger
	osClient    *opensearch.Client
	indexPrefix string
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			var event models.CDCEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				h.logger.Error("Failed to unmarshal event", zap.Error(err))
				session.MarkMessage(message, "")
				continue
			}

			// Process the event and index to OpenSearch
			if err := h.processEvent(event); err != nil {
				h.logger.Error("Failed to process event", zap.Error(err))
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *ConsumerGroupHandler) processEvent(event models.CDCEvent) error {
	// Extract entity type from the key
	parts := strings.Split(event.After.Key, "/")
	if len(parts) < 4 {
		h.logger.Error("Key has insufficient parts", zap.String("key", event.After.Key))
		return nil
	}
	entityType := parts[len(parts)-2] // e.g., "service", "node", "upstream"

	// Create index name based on entity type
	indexName := h.indexPrefix + "-" + entityType
	h.logger.Info("Processing event",
		zap.String("entityType", entityType),
		zap.String("indexName", indexName),
		zap.String("key", event.After.Key),
	)

	// Convert object to JSON
	objectBytes, err := json.Marshal(event.After.Value.Object)
	if err != nil {
		h.logger.Error("Failed to marshal object", zap.Error(err))
		return err
	}

	// Extract the ID from the object map
	obj, ok := event.After.Value.Object.(map[string]interface{})
	if !ok {
		h.logger.Error("Failed to cast object to map", zap.String("key", event.After.Key))
		return nil
	}

	id, ok := obj["id"].(string)
	if !ok {
		h.logger.Error("Failed to get ID from object", zap.String("key", event.After.Key))
		return nil
	}

	// Index the document
	h.logger.Info("Indexing document",
		zap.String("indexName", indexName),
		zap.String("id", id),
	)
	_, err = h.osClient.Index(
		indexName,
		strings.NewReader(string(objectBytes)),
		h.osClient.Index.WithDocumentID(id),
	)
	if err != nil {
		h.logger.Error("Failed to index document", zap.Error(err))
		return err
	}

	h.logger.Info("Successfully indexed document",
		zap.String("indexName", indexName),
		zap.String("id", id),
	)

	return nil
}
