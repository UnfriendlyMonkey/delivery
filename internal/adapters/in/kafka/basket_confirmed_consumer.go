package kafka

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/generated/queues/basketconfirmedpb"
	"delivery/internal/pkg/errs"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type BasketConfirmedConsumer interface {
	Consume() error
	Close() error
}

var (
	_ BasketConfirmedConsumer     = &basketConfirmedConsumer{}
	_ sarama.ConsumerGroupHandler = &basketConfirmedConsumer{}
)

type basketConfirmedConsumer struct {
	topic              string
	consumerGroup      sarama.ConsumerGroup
	ctx                context.Context
	cancel             context.CancelFunc
	createOrderHandler commands.CreateOrderHandler
}

func NewBasketConfirmedConsumer(
	brokers []string,
	group string,
	topic string,
	createOrderHandler commands.CreateOrderHandler,
) (BasketConfirmedConsumer, error) {
	if len(brokers) == 0 {
		return nil, errs.NewValueIsRequiredError("brokers")
	}
	if group == "" {
		return nil, errs.NewValueIsRequiredError("group")
	}
	if topic == "" {
		return nil, errs.NewValueIsRequiredError("topic")
	}
	if createOrderHandler == nil {
		return nil, errs.NewValueIsRequiredError("createOrderHandler")
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V3_4_0_0
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &basketConfirmedConsumer{
		topic: topic,
		consumerGroup: consumerGroup,
		createOrderHandler: createOrderHandler,
		ctx: ctx,
		cancel: cancel,
	}, nil
}

func (c *basketConfirmedConsumer) Close() error {
	c.cancel()
	return c.consumerGroup.Close()
}

func (c *basketConfirmedConsumer) Consume() error {
	for {
		err := c.consumerGroup.Consume(c.ctx, []string{c.topic}, c)
		if err != nil {
			log.Printf("Error consuming Kafka: %v", err)
			return err
		}
		if c.ctx.Err() != nil {
			return nil
		}
	}
}

// sarama.ConsumerGroupHandler interface implementation
func (c *basketConfirmedConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *basketConfirmedConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c *basketConfirmedConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		ctx := context.Background()
		fmt.Printf("Message received: topic = %s, partition = %d, offset = %d, key = %s, value = %s\n",
			message.Topic, message.Partition, message.Offset, message.Key, message.Value)

		var event basketconfirmedpb.BasketConfirmedIntegrationEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		cmd, err := commands.NewCreateOrderCommand(
			uuid.MustParse(event.BasketId), event.Address.Street, kernel.Volume(event.Volume),
		)
		if err != nil {
			log.Printf("Error creating createOrder command: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		if err := c.createOrderHandler.Handle(ctx, cmd); err != nil {
			log.Printf("Error handling createOrder command: %v", err)
		}

		session.MarkMessage(message, "")
	}

	return nil
}
