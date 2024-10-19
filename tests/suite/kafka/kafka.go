package kafka_suite

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/suite"
	kafka_client "gitlab.ozon.dev/chppppr/homework/internal/clients/kafka"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/infra/kafka"
	"gitlab.ozon.dev/chppppr/homework/internal/infra/kafka/consumer"
	"gitlab.ozon.dev/chppppr/homework/internal/infra/kafka/producer"
)

type KafkaSuite struct {
	suite.Suite
	pr        sarama.SyncProducer
	pr_client *kafka_client.ProducerClient
	cons      *consumer.Consumer

	result chan []byte
}

func (s *KafkaSuite) SetupSuite() {
	var err error
	topic := "pvz.events-log"
	kafka_cfg := kafka.Config{
		Brokers: []string{"localhost:9093"},
	}
	s.result = make(chan []byte, 1)

	s.pr, err = producer.NewSyncProducer(kafka_cfg)
	s.Require().NoError(err)
	s.pr_client = kafka_client.NewProducerClient(s.pr, topic)

	s.cons, err = consumer.NewConsumer(kafka_cfg)
	s.Require().NoError(err)

	err = s.cons.ConsumeTopic(context.Background(), topic, func(cm *sarama.ConsumerMessage) {
		s.result <- cm.Value
	}, &sync.WaitGroup{})
	s.Require().NoError(err)
}

func (s *KafkaSuite) TearDownSuite() {
	s.cons.Close()
	s.pr.Close()
}

func (s *KafkaSuite) TestEventOrderAccepted() {
	var err_user, err_ser error
	orders := []uint64{1}
	event_type := domain.EventOrderAccepted

	expected_event := domain.NewEvent(orders, event_type, err_user, err_ser)
	err := s.pr_client.Send(orders, event_type, err_user, err_ser)
	s.Require().NoError(err)

	actual_bytes := <-s.result
	var actual_event *domain.Event
	err = json.Unmarshal(actual_bytes, &actual_event)
	s.Require().NoError(err)

	s.Require().Equal(expected_event.OrderIDs, actual_event.OrderIDs)
	s.Require().Equal(expected_event.EventType, actual_event.EventType)
	s.Require().Equal(expected_event.ErrUser, actual_event.ErrUser)
	s.Require().Equal(expected_event.ErrService, actual_event.ErrService)
}

func (s *KafkaSuite) TestEventOrderGiveClient() {
	var err_user, err_ser error
	orders := []uint64{2}
	event_type := domain.EventOrderGiveClient

	expected_event := domain.NewEvent(orders, event_type, err_user, err_ser)
	err := s.pr_client.Send(orders, event_type, err_user, err_ser)
	s.Require().NoError(err)

	actual_bytes := <-s.result
	var actual_event *domain.Event
	err = json.Unmarshal(actual_bytes, &actual_event)
	s.Require().NoError(err)

	s.Require().Equal(expected_event.OrderIDs, actual_event.OrderIDs)
	s.Require().Equal(expected_event.EventType, actual_event.EventType)
	s.Require().Equal(expected_event.ErrUser, actual_event.ErrUser)
	s.Require().Equal(expected_event.ErrService, actual_event.ErrService)
}

func (s *KafkaSuite) TestEventOrderGiveReturned() {
	var err_user, err_ser error
	orders := []uint64{3}
	event_type := domain.EventOrderReturned

	expected_event := domain.NewEvent(orders, event_type, err_user, err_ser)
	err := s.pr_client.Send(orders, event_type, err_user, err_ser)
	s.Require().NoError(err)

	actual_bytes := <-s.result
	var actual_event *domain.Event
	err = json.Unmarshal(actual_bytes, &actual_event)
	s.Require().NoError(err)

	s.Require().Equal(expected_event.OrderIDs, actual_event.OrderIDs)
	s.Require().Equal(expected_event.EventType, actual_event.EventType)
	s.Require().Equal(expected_event.ErrUser, actual_event.ErrUser)
	s.Require().Equal(expected_event.ErrService, actual_event.ErrService)
}

func (s *KafkaSuite) TestEventOrderGiveCourier() {
	var err_user, err_ser error
	orders := []uint64{41, 42, 43, 44}
	event_type := domain.EventOrderGiveCourier

	expected_event := domain.NewEvent(orders, event_type, err_user, err_ser)
	err := s.pr_client.Send(orders, event_type, err_user, err_ser)
	s.Require().NoError(err)

	actual_bytes := <-s.result
	var actual_event *domain.Event
	err = json.Unmarshal(actual_bytes, &actual_event)
	s.Require().NoError(err)

	s.Require().Equal(expected_event.OrderIDs, actual_event.OrderIDs)
	s.Require().Equal(expected_event.EventType, actual_event.EventType)
	s.Require().Equal(expected_event.ErrUser, actual_event.ErrUser)
	s.Require().Equal(expected_event.ErrService, actual_event.ErrService)
}

func (s *KafkaSuite) TestEventOrderAcceptedWithUserError() {
	var err_user, err_ser error
	orders := []uint64{5}
	event_type := domain.EventOrderReturned
	err_user = domain.ErrNotFound

	expected_event := domain.NewEvent(orders, event_type, err_user, err_ser)
	err := s.pr_client.Send(orders, event_type, err_user, err_ser)
	s.Require().NoError(err)

	actual_bytes := <-s.result
	var actual_event *domain.Event
	err = json.Unmarshal(actual_bytes, &actual_event)
	s.Require().NoError(err)

	s.Require().Equal(expected_event.OrderIDs, actual_event.OrderIDs)
	s.Require().Equal(expected_event.EventType, actual_event.EventType)
	s.Require().Equal(expected_event.ErrUser, actual_event.ErrUser)
	s.Require().Equal(expected_event.ErrService, actual_event.ErrService)
}

func (s *KafkaSuite) TestEventOrderAcceptedWithServiceError() {
	var err_user, err_ser error
	orders := []uint64{6}
	event_type := domain.EventOrderReturned
	err_ser = fmt.Errorf("some service error")

	expected_event := domain.NewEvent(orders, event_type, err_user, err_ser)
	err := s.pr_client.Send(orders, event_type, err_user, err_ser)
	s.Require().NoError(err)

	actual_bytes := <-s.result
	var actual_event *domain.Event
	err = json.Unmarshal(actual_bytes, &actual_event)
	s.Require().NoError(err)

	s.Require().Equal(expected_event.OrderIDs, actual_event.OrderIDs)
	s.Require().Equal(expected_event.EventType, actual_event.EventType)
	s.Require().Equal(expected_event.ErrUser, actual_event.ErrUser)
	s.Require().Equal(expected_event.ErrService, actual_event.ErrService)
}
