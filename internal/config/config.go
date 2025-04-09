package config

// Config holds all configuration for the application
type Config struct {
	Kafka struct {
		Brokers  []string `mapstructure:"brokers"`
		Topic    string   `mapstructure:"topic"`
		GroupID  string   `mapstructure:"group_id"`
		ClientID string   `mapstructure:"client_id"`
		TopicConfig struct {
			Partitions        int `mapstructure:"partitions"`
			ReplicationFactor int `mapstructure:"replication_factor"`
		} `mapstructure:"topic_config"`
	} `mapstructure:"kafka"`

	OpenSearch struct {
		Hosts      []string `mapstructure:"hosts"`
		IndexPrefix string   `mapstructure:"index_prefix"`
	} `mapstructure:"opensearch"`

	Producer struct {
		InputFile string `mapstructure:"input_file"`
	} `mapstructure:"producer"`

	Consumer struct {
		BatchSize      int    `mapstructure:"batch_size"`
		CommitInterval string `mapstructure:"commit_interval"`
	} `mapstructure:"consumer"`

	Log struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
	} `mapstructure:"log"`
}
