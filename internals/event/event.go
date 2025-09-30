package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id             string `json:"id"`
	SessionId      string `json:"sessionId"`
	Type           string `json:"type"`
	SerializedData string `json:"serializedData"`
	CreatedAt      int64  `json:"createdAt"`
}

type EventKind string

const (
	CustomEvent EventKind = "custom"
)

func NewEvent(kind EventKind, data any) (Event, error) {
	serialized, err := json.Marshal(data)
	if err != nil {
		return Event{}, err
	}

	return Event{
		Id:             uuid.NewString(),
		SessionId:      "", // TODO: get current session id from session manager
		Type:           string(kind),
		SerializedData: string(serialized),
		CreatedAt:      time.Now().UnixMilli(),
	}, nil
}
