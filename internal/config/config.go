package config

import (
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Kafka      KafkaConfig      `mapstructure:"kafka"`
	OpenSearch OpenSearchConfig `mapstructure:"opensearch"`
	Input      InputConfig      `mapstructure:"input"`
}

// KafkaConfig holds Kafka-specific configuration
type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
}

// OpenSearchConfig holds OpenSearch-specific configuration
type OpenSearchConfig struct {
	Hosts       []string `mapstructure:"hosts"`
	IndexPrefix string   `mapstructure:"index_prefix"`
}

// InputConfig holds input file configuration
type InputConfig struct {
	FilePath string `mapstructure:"file_path"`
}

// LoadConfig loads the configuration from environment variables and config file
func LoadConfig() (*Config, error) {
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.topic", "cdc-events")
	viper.SetDefault("opensearch.hosts", []string{"http://localhost:9200"})
	viper.SetDefault("opensearch.index_prefix", "konnect")
	viper.SetDefault("input.file_path", "stream.jsonl")

	viper.AutomaticEnv()
	
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
