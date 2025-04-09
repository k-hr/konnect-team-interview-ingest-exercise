package config

import (
	"github.com/spf13/viper"
	"strings"
)

// LoadConfig loads the configuration from environment variables and config file
func LoadConfig() (*Config, error) {
	v := viper.New()

	// Set default configuration file path
	v.SetConfigName("application")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// Set defaults
	setDefaults(v)

	// Allow overrides from environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")

	// Handle environment variable overrides for arrays
	handleArrayOverrides(v)

	// Read configuration
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults sets the default values for configuration
func setDefaults(v *viper.Viper) {
	v.SetDefault("kafka.brokers", []string{"localhost:9092"})
	v.SetDefault("kafka.topic", "cdc-events")
	v.SetDefault("kafka.group_id", "cdc-consumer-group")
	v.SetDefault("kafka.client_id", "cdc-client")
	v.SetDefault("opensearch.hosts", []string{"http://localhost:9200"})
	v.SetDefault("opensearch.index_prefix", "cdc")
	v.SetDefault("producer.input_file", "stream.jsonl")
	v.SetDefault("consumer.batch_size", 100)
	v.SetDefault("consumer.commit_interval", "1s")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
}

// handleArrayOverrides handles environment variable overrides for array values
func handleArrayOverrides(v *viper.Viper) {
	if brokers := v.GetString("kafka.brokers"); brokers != "" {
		v.Set("kafka.brokers", strings.Split(brokers, ","))
	}
	if hosts := v.GetString("opensearch.hosts"); hosts != "" {
		v.Set("opensearch.hosts", strings.Split(hosts, ","))
	}
}
