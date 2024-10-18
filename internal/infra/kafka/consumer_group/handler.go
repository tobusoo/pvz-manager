package consumer_group

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

type ConsumerGroupHandler struct {
	logger *slog.Logger
}

func NewConsumerGroupHandler(logger *slog.Logger) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{logger: logger}
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) handleMessage(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) {
	var event domain.Event
	err := json.Unmarshal(message.Value, &event)
	if err != nil {
		h.logger.Error(fmt.Sprintf("can't unmarshal json: %v", err))
		return
	}

	if event.ErrService != "" {
		h.logger.Error(
			"handle message",
			"partition", message.Partition,
			"offset", message.Offset,
			"event", event.EventType,
			"timestamp", event.Timestamp,
			"orders_id", event.OrderIDs,
			"error_user", event.ErrUser,
			"error_serice", event.ErrService,
		)
	} else {
		h.logger.Info(
			"handle message",
			"partition", message.Partition,
			"offset", message.Offset,
			"event", event.EventType,
			"timestamp", event.Timestamp,
			"orders_id", event.OrderIDs,
			"error_user", event.ErrUser,
			"error_serice", event.ErrService,
		)
	}

	session.MarkMessage(message, "")
	session.Commit()
}

//gocognit:ignore
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			h.handleMessage(session, message)
		case <-session.Context().Done():
			return nil
		}
	}
}
