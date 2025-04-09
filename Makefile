.PHONY: build test clean run-producer run-consumer

build:
	go build -o bin/producer cmd/producer/main.go
	go build -o bin/consumer cmd/consumer/main.go

test:
	go test ./...

clean:
	rm -rf bin/

run-producer:
	go run cmd/producer/main.go

run-consumer:
	go run cmd/consumer/main.go

deps:
	go mod download

lint:
	golangci-lint run
