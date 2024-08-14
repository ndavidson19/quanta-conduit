package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const (
	maxRetries        = 5
	initialRetryDelay = 5 * time.Second
)

func (c *Consumer) ConsumeWithRetry(ctx context.Context) {
	for {
		if err := c.consume(ctx); err != nil {
			c.logger.Error("Kafka consumer error", zap.Error(err))
			if err := c.reconnect(ctx); err != nil {
				c.logger.Error("Failed to reconnect to Kafka", zap.Error(err))
				return
			}
		}

		if ctx.Err() != nil {
			return
		}
	}
}

func (c *Consumer) consume(ctx context.Context) error {
	for {
		if err := c.consumer.Consume(ctx, c.topics, c); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *Consumer) reconnect(ctx context.Context) error {
	retryDelay := initialRetryDelay

	for i := 0; i < maxRetries; i++ {
		c.logger.Info("Attempting to reconnect to Kafka", zap.Int("attempt", i+1))

		if err := c.consumer.Close(); err != nil {
			c.logger.Error("Failed to close existing consumer", zap.Error(err))
		}

		newConsumer, err := sarama.NewConsumerGroup(c.brokers, c.groupID, c.config)
		if err == nil {
			c.consumer = newConsumer
			c.logger.Info("Successfully reconnected to Kafka")
			return nil
		}

		c.logger.Error("Failed to reconnect to Kafka", zap.Error(err))

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(retryDelay):
			retryDelay *= 2
		}
	}

	return ErrMaxRetriesExceeded
}
