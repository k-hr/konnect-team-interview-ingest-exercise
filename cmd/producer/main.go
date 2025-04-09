package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/kong/konnect-ingest/internal/config"
	"github.com/kong/konnect-ingest/internal/models"
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

	// Create Kafka producer
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, kafkaConfig)
	if err != nil {
		logger.Fatal("Failed to create Kafka producer", zap.Error(err))
	}
	defer producer.Close()

	// Open input file
	file, err := os.Open(cfg.Input.FilePath)
	if err != nil {
		logger.Fatal("Failed to open input file", zap.Error(err))
	}
	defer file.Close()

	// Create decoder for JSON lines
	decoder := json.NewDecoder(file)

	// Setup signal handling for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Process events
	for decoder.More() {
		select {
		case <-signals:
			logger.Info("Received shutdown signal, stopping...")
			return
		default:
			var event models.CDCEvent
			if err := decoder.Decode(&event); err != nil {
				logger.Error("Failed to decode event", zap.Error(err))
				continue
			}

			// Convert event to JSON
			eventBytes, err := json.Marshal(event)
			if err != nil {
				logger.Error("Failed to marshal event", zap.Error(err))
				continue
			}

			// Send to Kafka
			msg := &sarama.ProducerMessage{
				Topic: cfg.Kafka.Topic,
				Value: sarama.StringEncoder(eventBytes),
			}

			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				logger.Error("Failed to send message", zap.Error(err))
				continue
			}

			logger.Info("Message sent",
				zap.Int32("partition", partition),
				zap.Int64("offset", offset),
			)
		}
	}
}
