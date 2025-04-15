package events

import (
	"github.com/google/uuid"
	"time"
)

type EventType string

const (
	LoginAttempt    EventType = "LoginAttempt"
	SuccessfulLogin EventType = "SuccessfulLogin"
	FailedLogin     EventType = "FailedLogin"
	Exit            EventType = "Exit"
	Blocking        EventType = "Blocking"
)

type Event struct {
	ID        string                 `json:"ID"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data"`
}

func NewEvent(eventType EventType, source string, data map[string]interface{}) Event {
	return Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Timestamp: time.Now(),
		Source:    source,
		Data:      data,
	}
}
