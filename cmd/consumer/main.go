package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/kong/konnect-ingest/internal/config"
	"github.com/kong/konnect-ingest/internal/consumer"
	"github.com/kong/konnect-ingest/internal/data_processing"
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

	kafkaConsumer, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.GroupID, kafkaConfig)
	if err != nil {
		logger.Fatal("Failed to create consumer group", zap.Error(err))
	}
	defer kafkaConsumer.Close()

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		logger.Info("Received shutdown signal, stopping...")
		cancel()
	}()

	// Create components
	indexer := data_processing.NewOpenSearchIndexer(osClient, logger)
	entityExtractor := data_processing.NewCDCEntityExtractor(logger)

	// Create event processor
	eventProcessor := consumer.NewCDCEventProcessor(
		logger,
		indexer,
		entityExtractor,
		cfg.OpenSearch.IndexPrefix,
	)

	// Create consumer handler
	consumerHandler := consumer.NewKafkaConsumerHandler(logger)
	consumerHandler.SetEventProcessor(eventProcessor)

	// Start consuming
	for {
		err := kafkaConsumer.Consume(ctx, []string{cfg.Kafka.Topic}, consumerHandler)
		if err != nil {
			if err == context.Canceled {
				return
			}
			logger.Error("Error from consumer", zap.Error(err))
		}
	}
}


