package broker

import amqp "github.com/rabbitmq/amqp091-go"

type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewRabbitMQ(conn *amqp.Connection) *RabbitMQ {
	ch, _ := conn.Channel()
	return &RabbitMQ{
		Conn: conn,
		Ch:   ch,
	}
}

func (r *RabbitMQ) Close() {
	if r.Ch != nil {
		r.Ch.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
}

func ConnectRabbit(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
