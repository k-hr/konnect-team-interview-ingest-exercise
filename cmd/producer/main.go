package main

import (
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/kong/konnect-ingest/internal/config"
	"github.com/kong/konnect-ingest/internal/producer"
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

	// Create event producer
	eventProducer, err := producer.NewKafkaEventProducer(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		logger,
	)
	if err != nil {
		logger.Fatal("Failed to create event producer", zap.Error(err))
	}
	defer eventProducer.Close()

	// Create event reader
	eventReader, err := producer.NewEventReader(cfg.Producer.InputFile, logger)
	if err != nil {
		logger.Fatal("Failed to create event reader", zap.Error(err))
	}
	defer eventReader.Close()

	// Setup signal handling for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Process events
	for {
		select {
		case <-signals:
			logger.Info("Received shutdown signal, stopping...")
			return
		default:
			event, err := eventReader.ReadEvent()
			if err == io.EOF {
				return
			}
			if err != nil {
				logger.Error("Failed to read event", zap.Error(err))
				continue
			}

			if err := eventProducer.ProduceEvent(*event); err != nil {
				logger.Error("Failed to produce event", zap.Error(err))
				continue
			}
		}
	}
}
