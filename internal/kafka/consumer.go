package kafka

import (
	"context"
	"encoding/json"

	"conduit/internal/websocket"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Consumer struct {
	consumer  sarama.ConsumerGroup
	topics    []string
	logger    *zap.Logger
	wsManager *websocket.Manager
}

func NewConsumer(brokers []string, groupID string, topics []string, logger *zap.Logger, wsManager *websocket.Manager) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:  consumer,
		topics:    topics,
		logger:    logger,
		wsManager: wsManager,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context) {
	for {
		if err := c.consumer.Consume(ctx, c.topics, c); err != nil {
			c.logger.Error("Error from consumer", zap.Error(err))
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		c.logger.Info("Message claimed",
			zap.String("topic", message.Topic),
			zap.Int32("partition", message.Partition),
			zap.Int64("offset", message.Offset),
		)

		// Process the message and broadcast to WebSocket clients
		var data map[string]interface{}
		if err := json.Unmarshal(message.Value, &data); err != nil {
			c.logger.Error("Failed to unmarshal message", zap.Error(err))
			continue
		}

		broadcastMessage, err := json.Marshal(map[string]interface{}{
			"type": "market_data",
			"data": data,
		})
		if err != nil {
			c.logger.Error("Failed to marshal broadcast message", zap.Error(err))
			continue
		}

		c.wsManager.broadcast <- broadcastMessage

		session.MarkMessage(message, "")
	}
	return nil
}
