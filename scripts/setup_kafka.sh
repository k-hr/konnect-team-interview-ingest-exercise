#!/bin/bash

# Read configuration from application.yml
if [ ! -f "application.yml" ]; then
    echo "Error: application.yml not found"
    exit 1
fi

# Extract values using yq (assumes yq is installed)
TOPIC=$(yq eval '.kafka.topic' application.yml)
PARTITIONS=$(yq eval '.kafka.topic_config.partitions' application.yml)
REPLICATION_FACTOR=$(yq eval '.kafka.topic_config.replication_factor' application.yml)
# Get broker from config but use internal Docker network address
BOOTSTRAP_SERVER="kafka:29092"

# Validate configuration
if [ -z "$TOPIC" ] || [ -z "$PARTITIONS" ] || [ -z "$REPLICATION_FACTOR" ] || [ -z "$BOOTSTRAP_SERVER" ]; then
    echo "Error: Missing required configuration in application.yml"
    exit 1
fi

# Create the topic if it doesn't exist
docker exec konnect-team-interview-ingest-exercise-kafka-1 kafka-topics --create \
  --if-not-exists \
  --topic $TOPIC \
  --bootstrap-server $BOOTSTRAP_SERVER \
  --partitions $PARTITIONS \
  --replication-factor $REPLICATION_FACTOR

# Describe the topic to verify
docker exec konnect-team-interview-ingest-exercise-kafka-1 kafka-topics --describe \
  --topic $TOPIC \
  --bootstrap-server $BOOTSTRAP_SERVER
