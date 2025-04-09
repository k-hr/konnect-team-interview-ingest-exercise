.PHONY: all build test clean run-producer run-consumer init config setup-kafka init-all

# Default target
all: config build

# Full initialization with Kafka setup
init-all: init setup-kafka

# Configuration paths
CONFIG_SAMPLE := application.yml.sample
CONFIG_FILE := application.yml

# Build targets
build: config
	go build -o bin/producer cmd/producer/main.go
	go build -o bin/consumer cmd/consumer/main.go

# Initialize project
init: deps config

# Test target
test:
	go test ./...

# Clean target
clean:
	rm -rf bin/
	rm -f $(CONFIG_FILE)

# Configuration management
config:
	@if [ ! -f $(CONFIG_FILE) ]; then \
		echo "Creating config file from sample..."; \
		cp $(CONFIG_SAMPLE) $(CONFIG_FILE); \
	fi

# Run targets with config
run-producer: config
	CONFIG_FILE=$(CONFIG_FILE) ./bin/producer

run-consumer: config
	CONFIG_FILE=$(CONFIG_FILE) ./bin/consumer

# Dependency management
deps:
	go mod download
	go mod tidy

# Kafka setup
setup-kafka: config
	@command -v yq >/dev/null 2>&1 || { echo "yq is required but not installed. Install with: brew install yq"; exit 1; }
	@docker ps | grep -q konnect-team-interview-ingest-exercise-kafka-1 || { echo "Kafka container not running. Please start with: docker-compose up -d"; exit 1; }
	@chmod +x scripts/setup_kafka.sh
	./scripts/setup_kafka.sh

# Linting
lint:
	golangci-lint run
