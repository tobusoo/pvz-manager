package domain

import "time"

type EventType string

const (
	EventOrderAccepted    EventType = "order accepted"
	EventOrderGiveClient  EventType = "order issued to client"
	EventOrderGiveCourier EventType = "order issued to courier"
	EventOrderReturned    EventType = "order returned"
)

type Event struct {
	EventType EventType `json:"event"`
	Timestamp time.Time `json:"timestamp"`

	OrderIDs   []uint64 `json:"orders_id"`
	ErrUser    string   `json:"error_user"`
	ErrService string   `json:"error_service"`
}

func errToString(err error) string {
	if err != nil {
		return err.Error()
	}

	return ""
}

func NewEvent(orderIDs []uint64, eventType EventType, err_user, err_ser error) *Event {
	return &Event{
		EventType: eventType,
		Timestamp: time.Now().UTC(),

		OrderIDs:   orderIDs,
		ErrUser:    errToString(err_user),
		ErrService: errToString(err_ser),
	}
}
