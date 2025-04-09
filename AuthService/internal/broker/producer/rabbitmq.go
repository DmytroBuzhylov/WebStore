package producer

import (
	"AuthService/internal/broker"
	"context"
)

type Producer struct {
	rabbit *broker.RabbitMQ
}

func NewProducer(rabbit *broker.RabbitMQ) *Producer {
	return &Producer{rabbit: rabbit}
}

func (p *Producer) SendVerificationCode(ctx context.Context, email string, code string) error {
	return nil
}
