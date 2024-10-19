package consumer

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/chppppr/homework/internal/infra/kafka"
)

type Consumer struct {
	consumer sarama.Consumer
}

func NewConsumer(kafkaConfig kafka.Config) (*Consumer, error) {
	config := sarama.NewConfig()

	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(kafkaConfig.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
	}, err
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

//gocyclo:ignore
//gocognit:ignore
func (c *Consumer) ConsumeTopic(ctx context.Context, topic string, handler func(*sarama.ConsumerMessage), wg *sync.WaitGroup) error {
	partitionList, err := c.consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for _, partition := range partitionList {
		pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		wg.Add(1)
		go func(pc sarama.PartitionConsumer, partition int32) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-pc.Messages():
					if !ok {
						return
					}
					handler(msg)
				}
			}
		}(pc, partition)
	}

	return nil
}
