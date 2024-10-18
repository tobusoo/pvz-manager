package consumer_group

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type consumerGroup struct {
	sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
}

//gocognit:ignore
func (c *consumerGroup) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := c.ConsumerGroup.Consume(ctx, c.topics, c.handler); err != nil {
				log.Printf("ConsumerGroup.Consume(): %v", err)
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()
}

func NewConsumerGroup(brokers []string, groupID string, topics []string, consumerGroupHandler sarama.ConsumerGroupHandler, opts ...Option) (*consumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	config.Consumer.Group.ResetInvalidOffsets = true
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	config.Consumer.Return.Errors = true

	for _, opt := range opts {
		opt.Apply(config)
	}

	cg, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &consumerGroup{
		ConsumerGroup: cg,
		handler:       consumerGroupHandler,
		topics:        topics,
	}, nil
}
