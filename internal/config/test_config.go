package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Set up environment variables
	os.Setenv("APP_DATABASE_URL", "postgres://user:pass@localhost:5432/dbname")
	os.Setenv("APP_SERVER_PORT", "8080")
	os.Setenv("APP_JWT_SECRET", "test_secret")
	os.Setenv("APP_KAFKA_BROKERS", "localhost:9092,localhost:9093")
	os.Setenv("APP_KAFKA_GROUP_ID", "test_group")
	os.Setenv("APP_KAFKA_TOPICS", "topic1,topic2")

	// Test
	cfg, err := Load()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "postgres://user:pass@localhost:5432/dbname", cfg.DatabaseURL)
	assert.Equal(t, "8080", cfg.ServerPort)
	assert.Equal(t, "test_secret", cfg.JWTSecret)
	assert.Equal(t, []string{"localhost:9092", "localhost:9093"}, cfg.KafkaBrokers)
	assert.Equal(t, "test_group", cfg.KafkaGroupID)
	assert.Equal(t, []string{"topic1", "topic2"}, cfg.KafkaTopics)
}

func TestLoadMissingRequired(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	// Test
	_, err := Load()

	// Assert
	assert.Error(t, err)
}