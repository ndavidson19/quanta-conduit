package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	DatabaseURL  string   `envconfig:"DATABASE_URL" required:"true"`
	ServerPort   string   `envconfig:"SERVER_PORT" default:"8080"`
	JWTSecret    string   `envconfig:"JWT_SECRET" required:"true"`
	KafkaBrokers []string `envconfig:"KAFKA_BROKERS" required:"true"`
	KafkaGroupID string   `envconfig:"KAFKA_GROUP_ID" required:"true"`
	KafkaTopics  []string `envconfig:"KAFKA_TOPICS" required:"true"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("app", &cfg)
	return &cfg, err
}