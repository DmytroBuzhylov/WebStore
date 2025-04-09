package event

import "context"

type Publisher interface {
	Publish(ctx context.Context, msg EventMessage) error
}

type EventMessage struct {
	EventId        string `json:"eventId"`
	EventType      string `json:"eventType"`
	EventTimestamp string `json:"eventTimestamp"`
	Version        string `json:"version"`
	Payload        UserRegisterPayload
}

type UserRegisterPayload struct {
	UserID       string `json:"userID"`
	Email        string `json:"email"`
	FullName     string `json:"fullName"`
	RegisteredAt string `json:"registeredAt"`
}
