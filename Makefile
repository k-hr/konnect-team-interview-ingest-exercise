.PHONY: build test clean run-producer run-consumer init config

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

# Linting
lint:
	golangci-lint run
