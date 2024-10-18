package kafka_client

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

type ProducerClient struct {
	prod  sarama.SyncProducer
	topic string
}

func NewProducerClient(producer sarama.SyncProducer, topic string) *ProducerClient {
	return &ProducerClient{
		prod:  producer,
		topic: topic,
	}
}

func (p *ProducerClient) Send(orderIDs []uint64, eventType domain.EventType, err_usr, err_ser error) error {
	ev := domain.NewEvent(orderIDs, eventType, err_usr, err_ser)

	bytes, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Value:     sarama.ByteEncoder(bytes),
		Timestamp: time.Now(),
	}

	_, _, err = p.prod.SendMessage(msg)
	return err
}
