package rabbitmq

import (
    "fmt"
    "github.com/streadway/amqp"
)

type Publisher struct {
    connection *amqp.Connection
    channel    *amqp.Channel
}

// NewPublisher creates a new RabbitMQ publisher
func NewPublisher(amqpURL string) (*Publisher, error) {
    conn, err := amqp.Dial(amqpURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to open a channel: %v", err)
    }

    return &Publisher{
        connection: conn,
        channel:    ch,
    }, nil
}

// Publish sends a message to the specified exchange and routing key
func (p *Publisher) Publish(exchange, routingKey string, body []byte) error {
    err := p.channel.Publish(
        exchange,   // exchange
        routingKey, // routing key
        false,      // mandatory
        false,      // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        })
    if err != nil {
        return fmt.Errorf("failed to publish message: %v", err)
    }
    return nil
}

// SetupRabbitMQ configures the exchange and queue for the publisher
func (p *Publisher) SetupRabbitMQ() error {
    // Declare the exchange
    err := p.channel.ExchangeDeclare(
        "verification", // exchange name
        "direct",       // type of exchange
        true,           // durable
        false,          // auto-delete
        false,          // internal
        false,          // no-wait
        nil,            // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare exchange: %v", err)
    }

    // Declare the queue and bind it to the exchange with the routing key
    _, err = p.channel.QueueDeclare(
        "otp_queue", // queue name
        true,        // durable
        false,       // delete when unused
        false,       // exclusive
        false,       // no-wait
        nil,         // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare queue: %v", err)
    }

    err = p.channel.QueueBind(
        "otp_queue",       // queue name
        "otp_routing_key", // routing key
        "verification",    // exchange name
        false,
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed to bind queue to exchange: %v", err)
    }

    return nil
}




// Close closes the RabbitMQ connection and channel
func (p *Publisher) Close() {
    p.channel.Close()
    p.connection.Close()
}
